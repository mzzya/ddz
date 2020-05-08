package tracer

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v7"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

type redisHook struct {
	config RedisHookConfig
	dbType string
}

// RedisHookConfig .
type RedisHookConfig struct {
	Name string
}

// NewRedisHook .
func NewRedisHook(config RedisHookConfig) redis.Hook {
	return &redisHook{config: config, dbType: "redis"}
}

// BeforeProcess .
func (r *redisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	span, newCtx := opentracing.StartSpanFromContext(ctx, fmt.Sprintf("%s:%s", r.config.Name, cmd.Name()))
	if span == nil {
		span = opentracing.StartSpan(cmd.Name())
		newCtx = opentracing.ContextWithSpan(ctx, span)
	}
	ext.DBType.Set(span, r.dbType)
	ext.DBInstance.Set(span, r.config.Name)
	ext.DBStatement.Set(span, cmd.String())
	return newCtx, nil
}

// AfterProcess .
func (r *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	err := cmd.Err()
	if err != nil {
		ext.Error.Set(span, true)
		span.LogFields(log.Error(err))
	}
	span.Finish()
	return nil
}
func (r *redisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}
func (r *redisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}

// NewTracerRedisClient .
func NewTracerRedisClient(ctx context.Context, cli *redis.Client) *redis.Client {
	if !Enable {
		return cli
	}
	return cli.WithContext(ctx)
}

// NewTracerRedisClusterClient .
func NewTracerRedisClusterClient(ctx context.Context, cli *redis.ClusterClient) *redis.ClusterClient {
	if !Enable {
		return cli
	}
	return cli.WithContext(ctx)
}

// NewTracerRedisSentinelClient .
func NewTracerRedisSentinelClient(ctx context.Context, cli *redis.SentinelClient) *redis.SentinelClient {
	if !Enable {
		return cli
	}
	return cli.WithContext(ctx)
}
