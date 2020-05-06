package tracer

import (
	"testing"

	"github.com/go-redis/redis/v7"
	"github.com/stretchr/testify/assert"
)

var redisCli *redis.Client

func initTestRedis() {
	redisCli = redis.NewClient(&redis.Options{})
	status := redisCli.Ping()
	err := status.Err()
	if err != nil {
		panic(err)
	}
}

//docker run --restart always  -d  -p 6379:6379 --name redis redis:6.0.1-alpine
func Test_redisHook_BeforeProcess(t *testing.T) {
	initTestRedis()
	hook := NewRedisHook(RedisHookConfig{Name: "test_redis"})
	redisCli.AddHook(hook)
	res, err := redisCli.Info().Result()
	assert.Equal(t, nil, err)
	t.Logf("Info Result:%s", res)
}
