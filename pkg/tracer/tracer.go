package tracer

import (
	"os"

	"github.com/spf13/viper"
	"github.com/uber/jaeger-client-go/config"
)

var (
	//enable 是否启用Opentracing
	enable bool
)

const (
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

// Enable 启用链路跟踪
func Enable(v *viper.Viper) {
	enable = v.GetBool("OPENTRACING_ENABLE")
	if !enable {
		return
	}
	cfg, err := ConfigFromViper(v)
	if err != nil {

	}
}

// ConfigFromViper 从viper中获取配置信息
func ConfigFromViper(v *viper.Viper) (c *config.Configuration, err error) {
	os.Setenv(envServiceName, v.GetString(envServiceName))
	os.Setenv(envDisabled, v.GetString(envDisabled))
	os.Setenv(envRPCMetrics, v.GetString(envRPCMetrics))
	os.Setenv(envTags, v.GetString(envTags))
	os.Setenv(envSamplerType, v.GetString(envSamplerType))
	os.Setenv(envSamplerParam, v.GetString(envSamplerParam))
	os.Setenv(envSamplerManagerHostPort, v.GetString(envSamplerManagerHostPort))
	os.Setenv(envSamplerMaxOperations, v.GetString(envSamplerMaxOperations))
	os.Setenv(envSamplerRefreshInterval, v.GetString(envSamplerRefreshInterval))
	os.Setenv(envReporterMaxQueueSize, v.GetString(envReporterMaxQueueSize))
	os.Setenv(envReporterFlushInterval, v.GetString(envReporterFlushInterval))
	os.Setenv(envReporterLogSpans, v.GetString(envReporterLogSpans))
	os.Setenv(envEndpoint, v.GetString(envEndpoint))
	os.Setenv(envUser, v.GetString(envUser))
	os.Setenv(envPassword, v.GetString(envPassword))
	os.Setenv(envAgentHost, v.GetString(envAgentHost))
	os.Setenv(envAgentPort, v.GetString(envAgentPort))
	c, err = config.FromEnv()
	return
}
