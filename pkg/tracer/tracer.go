package tracer

import (
	"io"
	"os"

	"github.com/hellojqk/simple_api/pkg/logger"
	"github.com/hellojqk/simple_api/pkg/util"
	"github.com/opentracing/opentracing-go"
	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
)

var (
	//Enable 是否启用Opentracing
	Enable       bool
	tracer       opentracing.Tracer
	tracerCloser io.Closer
)

const (
	confOpentracingEnable = "OPENTRACING_ENABLE"
	// environment variable names
	envServiceName            = "JAEGER_SERVICE_NAME"
	envDisabled               = "JAEGER_DISABLED"
	envRPCMetrics             = "JAEGER_RPC_METRICS"
	envTags                   = "JAEGER_TAGS"
	envSamplerType            = "JAEGER_SAMPLER_TYPE"
	envSamplerParam           = "JAEGER_SAMPLER_PARAM"
	envSamplerManagerHostPort = "JAEGER_SAMPLER_MANAGER_HOST_PORT"
	envSamplerMaxOperations   = "JAEGER_SAMPLER_MAX_OPERATIONS"
	envSamplerRefreshInterval = "JAEGER_SAMPLER_REFRESH_INTERVAL"
	envReporterMaxQueueSize   = "JAEGER_REPORTER_MAX_QUEUE_SIZE"
	envReporterFlushInterval  = "JAEGER_REPORTER_FLUSH_INTERVAL"
	envReporterLogSpans       = "JAEGER_REPORTER_LOG_SPANS"
	envEndpoint               = "JAEGER_ENDPOINT"
	envUser                   = "JAEGER_USER"
	envPassword               = "JAEGER_PASSWORD"
	envAgentHost              = "JAEGER_AGENT_HOST"
	envAgentPort              = "JAEGER_AGENT_PORT"
)

// Init 启用链路跟踪
func Init(v *viper.Viper) {
	Enable = v.GetBool(confOpentracingEnable)
	if !Enable {
		return
	}
	cfg, err := ConfigFromViper(v)
	if err != nil {
		logger.Logger.Error("tracer enable from viper", zap.Error(err))
	}
	options := make([]config.Option, 0, 3)
	os.Stdin.Chmod(os.ModeSetuid)
	os.Stdout.Chmod(os.ModeSetuid)
	options = append(options, config.Logger(jaeger.StdLogger))
	options = append(options, config.Metrics(prometheus.New()))
	util.PrintJSONWithColor(cfg)
	tracer, tracerCloser, err = cfg.NewTracer(options...)
	if err != nil {
		logger.Logger.Error("tracer enable create", zap.Error(err))
	}
	opentracing.SetGlobalTracer(tracer)
	util.CloserAdd(99, Closer)
}

// ConfigFromViper 从viper中获取配置信息
func ConfigFromViper(v *viper.Viper) (c *config.Configuration, err error) {
	setEnv(envServiceName, v.GetString(envServiceName))
	setEnv(envDisabled, v.GetString(envDisabled))
	setEnv(envRPCMetrics, v.GetString(envRPCMetrics))
	setEnv(envTags, v.GetString(envTags))
	setEnv(envSamplerType, v.GetString(envSamplerType))
	setEnv(envSamplerParam, v.GetString(envSamplerParam))
	setEnv(envSamplerManagerHostPort, v.GetString(envSamplerManagerHostPort))
	setEnv(envSamplerMaxOperations, v.GetString(envSamplerMaxOperations))
	setEnv(envSamplerRefreshInterval, v.GetString(envSamplerRefreshInterval))
	setEnv(envReporterMaxQueueSize, v.GetString(envReporterMaxQueueSize))
	setEnv(envReporterFlushInterval, v.GetString(envReporterFlushInterval))
	setEnv(envReporterLogSpans, v.GetString(envReporterLogSpans))
	setEnv(envEndpoint, v.GetString(envEndpoint))
	setEnv(envUser, v.GetString(envUser))
	setEnv(envPassword, v.GetString(envPassword))
	setEnv(envAgentHost, v.GetString(envAgentHost))
	setEnv(envAgentPort, v.GetString(envAgentPort))
	c, err = config.FromEnv()
	return
}

//setEnv .
func setEnv(key string, value string) error {
	err := os.Setenv(key, value)
	if err != nil {
		logger.Logger.Error("setenv", zap.Error(err))
	}
	return nil
}

// Closer .
func Closer() (errs []error) {
	errs = make([]error, 0)
	if err := tracerCloser.Close(); err != nil {
		errs = append(errs, err)
	}
	return
}
