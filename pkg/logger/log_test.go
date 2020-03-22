package logger

import (
	"testing"

	"github.com/spf13/viper"
)

func TestMain(m *testing.M) {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}
	Init(v)
	m.Run()
}

func TestInit(t *testing.T) {
	Logger.Debug("Debug log")
	Logger.Info("Info log")
	Logger.Warn("Warn log")
	Logger.Error("Error log")
	Logger.DPanic("DPanic log")
	// Logger.Panic("Panic log")
	// Logger.Fatal("Fatal log")
	// Logger.Sync()
}

// go test -bench . -run Logger  -benchmem
func BenchmarkLogger(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Logger.Debug("Debug log")
		Logger.Info("Info log")
		Logger.Warn("Warn log")
		Logger.Error("Error log")
		Logger.DPanic("DPanic log")
	}
}
