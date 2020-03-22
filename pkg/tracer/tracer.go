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
	os.Set(envServiceName, v.GetString(envServiceName))
	os.Set(envDisabled, v.GetString(envDisabled))
	os.Set(envRPCMetrics, v.GetString(envRPCMetrics))
	os.Set(envTags, v.GetString(envTags))
	os.Set(envSamplerType, v.GetString(envSamplerType))
	os.Set(envSamplerParam, v.GetString(envSamplerParam))
	os.Set(envSamplerManagerHostPort, v.GetString(envSamplerManagerHostPort))
	os.Set(envSamplerMaxOperations, v.GetString(envSamplerMaxOperations))
	os.Set(envSamplerRefreshInterval, v.GetString(envSamplerRefreshInterval))
	os.Set(envReporterMaxQueueSize, v.GetString(envReporterMaxQueueSize))
	os.Set(envReporterFlushInterval, v.GetString(envReporterFlushInterval))
	os.Set(envReporterLogSpans, v.GetString(envReporterLogSpans))
	os.Set(envEndpoint, v.GetString(envEndpoint))
	os.Set(envUser, v.GetString(envUser))
	os.Set(envPassword, v.GetString(envPassword))
	os.Set(envAgentHost, v.GetString(envAgentHost))
	os.Set(envAgentPort, v.GetString(envAgentPort))
	c, err = config.FromEnv()
	return
}
