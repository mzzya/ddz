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
	span, newCtx := opentracing.StartSpanFromContext(ctx, r.config.Name)
	if span == nil {
		span = opentracing.StartSpan(cmd.Name())
		newCtx = opentracing.ContextWithSpan(ctx, span)
	}
	ext.DBType.Set(span, r.dbType)
	ext.DBInstance.Set(span, cmd.Name())
	fmt.Printf("%v", cmd.Args())
	ext.DBStatement.Set(span, cmd.Name())
	return newCtx, nil
}

// AfterProcess .
func (r *redisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return nil
	}
	err := cmd.Err()
	if err != nil && err != redis.Nil {
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
