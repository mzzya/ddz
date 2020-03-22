package logger

import (
	"os"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/hellojqk/simple_api/pkg/util"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	//enableFileLogger 启用文件记录器
	enableFileLogger bool
	//enableKafkaLogger 启用kafka记录器
	enableKafkaLogger bool
	//enableConsoleLogger 启用控制台记录器
	enableConsoleLogger bool
	//Logger 日志记录器
	Logger      *zap.Logger
	fileLogger  *lumberjack.Logger
	kafkaLogger *KafkaLogger
)

const (
	confEnableFileLogger     = "LOGGER_ENABLE_FILE_LOGGER"
	confEnableKafkaLogger    = "LOGGER_ENABLE_KAFKA_LOGGER"
	confEnableConsoleLogger  = "LOGGER_ENABLE_CONSOLE_LOGGER"
	confLoggerLevel          = "LOGGER_LEVEL"
	confLoggerFileName       = "LOGGER_FILE_NAME"
	confLoggerFileMaxSize    = "LOGGER_FILE_MAX_SIZE"
	confLoggerFileMaxAge     = "LOGGER_FILE_MAX_AGE"
	confLoggerFileMaxBackups = "LOGGER_FILE_MAX_BACKUPS"
	confLoggerFileLocalTime  = "LOGGER_FILE_LOCAL_TIME"
	confLoggerFileCompress   = "LOGGER_FILE_COMPRESS"

	confLoggerKafkaAddress = "LOGGER_KAFKA_ADDRESS"
	confLoggerKafkaTopic   = "LOGGER_KAFKA_TOPIC"
)

// Init 日志初始化
func Init(v *viper.Viper) {
	cores := make([]zapcore.Core, 0, 3)

	//日志记录级别
	loggerLevel := zapcore.Level(v.GetInt(confLoggerLevel))
	prdEncoderConfig := zap.NewProductionEncoderConfig()
	prdEncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	if enableFileLogger = v.GetBool(confEnableFileLogger); enableFileLogger {
		fileLogger = &lumberjack.Logger{
			Filename:   v.GetString(confLoggerFileName),
			MaxSize:    v.GetInt(confLoggerFileMaxSize),
			MaxBackups: v.GetInt(confLoggerFileMaxBackups),
			MaxAge:     v.GetInt(confLoggerFileMaxAge),
			LocalTime:  v.GetBool(confLoggerFileLocalTime),
			Compress:   v.GetBool(confLoggerFileCompress),
		}
		encoder := zapcore.NewJSONEncoder(prdEncoderConfig)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(fileLogger), loggerLevel))
	}

	if enableKafkaLogger = v.GetBool(confEnableKafkaLogger); enableKafkaLogger {
		conf := sarama.NewConfig()
		address := v.GetStringSlice(confLoggerKafkaAddress)
		//等待服务器所有副本都保存成功后的响应
		conf.Producer.RequiredAcks = sarama.NoResponse
		//随机的分区类型
		conf.Producer.Partitioner = sarama.NewRandomPartitioner
		//是否等待成功和失败后的响应,只有上面的RequireAcks设置不是NoReponse这里才有用.
		conf.Producer.Return.Successes = true
		conf.Producer.Return.Errors = true
		producer, err := sarama.NewSyncProducer(address, conf)
		if err != nil {
			panic(errors.WithMessage(err, "logger init kafka producer"))
		}
		kafkaLogger = &KafkaLogger{
			Producer: producer,
			Topic:    v.GetString(confLoggerKafkaTopic),
		}
		encoder := zapcore.NewJSONEncoder(prdEncoderConfig)
		cores = append(cores, zapcore.NewCore(encoder, zapcore.AddSync(kafkaLogger), loggerLevel))
	}

	if enableConsoleLogger = v.GetBool(confEnableConsoleLogger); enableConsoleLogger {
		encoder := zapcore.NewConsoleEncoder(prdEncoderConfig)
		core := zapcore.NewCore(encoder, os.Stdout, loggerLevel)
		cores = append(cores, core)
	}
	opts := []zap.Option{}
	opts = append(opts, zap.AddCaller())
	opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	//采样率配置 1秒钟内重复msg达到100次时每隔10次记录一次
	opts = append(opts, zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewSampler(core, time.Second, 100, 10)
	}))
	Logger = zap.New(zapcore.NewTee(cores...)).WithOptions(opts...)
	// {"level":"info","ts":1584886046.821289,"caller":"log/log_test.go:37","msg":"Info log"}
	//  375200             10199 ns/op              22 B/op          0 allocs/op 	4s
	// Logger, _ = zap.NewProduction()
	util.Add(100, Close)
}

// KafkaLogger .
type KafkaLogger struct {
	Producer sarama.SyncProducer
	Topic    string
}

var kafkaMsgPool = sync.Pool{
	New: func() interface{} {
		return &sarama.ProducerMessage{}
	},
}

func (m *KafkaLogger) Write(p []byte) (n int, err error) {
	gotMsg := kafkaMsgPool.Get()
	msg := gotMsg.(*sarama.ProducerMessage)
	// msg := &sarama.ProducerMessage{}
	msg.Topic = m.Topic
	msg.Value = sarama.ByteEncoder(p)
	_, _, err = m.Producer.SendMessage(msg)
	kafkaMsgPool.Put(msg)
	if err != nil {
		return
	}
	return
}

// Close .
func Close() (errs []error) {
	errs = make([]error, 0, 2)
	if err := Logger.Sync(); err != nil {
		errs = append(errs, errors.WithMessage(err, "logger sync"))
	}
	if fileLogger != nil {
		if err := fileLogger.Close(); err != nil {
			errs = append(errs, errors.WithMessage(err, "logger file sync"))
		}
	}
	return
}
