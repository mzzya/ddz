package tracer

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	req, err := http.NewRequest("GET", "https://www.baidu.com", nil)
	assert.Equal(t, nil, err)
	span, err := DoBefore(context.Background(), req, "http_get")
	assert.Equal(t, nil, err)
	cli := http.Client{}
	resp, err := cli.Do(req)
	assert.Equal(t, nil, err)
	DoAfter(span, resp, err)
}
