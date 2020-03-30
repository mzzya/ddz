package cache

import (
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestRedis(t *testing.T) {
	res, err := redisCli.Info().Result()
	assert.Equal(t, nil, err)
	t.Logf("%s\n", res)
}
