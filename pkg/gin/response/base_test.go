package response

import (
	"context"
	"testing"

	"github.com/hellojqk/simple_api/pkg/code"
)

func TestRun(t *testing.T) {
	NewResponse(context.Background(), code.Success, nil)
}
