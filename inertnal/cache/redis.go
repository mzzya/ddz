package cache

import (
	"github.com/go-redis/redis/v7"
	"github.com/hellojqk/simple/pkg/logger"
	"go.uber.org/zap"
)

var redisCli *redis.Client

func init() {
	// config.DefaultViper()
	redisCli = redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	err := redisCli.Ping().Err()
	if err != nil {
		logger.Logger.Error("redis init", zap.Error(err))
	}
}
