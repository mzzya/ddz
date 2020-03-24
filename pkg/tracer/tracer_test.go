package tracer

import (
	"fmt"
	"testing"

	"github.com/hellojqk/simple_api/pkg/config"
	"github.com/hellojqk/simple_api/pkg/logger"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func TestMain(m *testing.M) {
	v := config.DefaultViper()
	logger.Init(v)
	Init(v)
	m.Run()
	// util.Close()
}

func TestInit(t *testing.T) {
	span := opentracing.StartSpan("tracer_test2")
	span.LogFields(log.Int("log1", 2222))
	span.Finish()
	fmt.Print("tracerCloser.Close().Error()", tracerCloser.Close())
}
