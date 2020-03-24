package logger

import (
	"testing"

	"github.com/hellojqk/simple_api/pkg/config"
	"github.com/hellojqk/simple_api/pkg/util"
)

func TestMain(m *testing.M) {
	Init(config.DefaultViper())
	m.Run()
	util.Close()
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
