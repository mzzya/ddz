package tracer

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDo(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	assert.Equal(t, nil, err)
	DoBefore()
}
