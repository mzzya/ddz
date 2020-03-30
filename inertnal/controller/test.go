package controller

import (
	"context"
	"sync"

	ginext "github.com/hellojqk/simple/pkg/gin-ext"
)

type TestResponse struct {
	ginext.BaseResponse
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &TestRequest{}
	},
}

type TestRequest struct {
	ginext.BaseRequest
}

func (r *TestRequest) New() ginext.Process {
	return &TestRequest{}
}
func (r *TestRequest) Exec(ctx context.Context) interface{} {
	resp := TestResponse{}
	resp.BaseResponse = ginext.NewSuccessResponse(ctx)
	return resp
}
