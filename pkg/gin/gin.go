package gin

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple_api/pkg/code"
	"github.com/hellojqk/simple_api/pkg/gin/response"
)

// Process 请求封装
type Process interface {
	New() Process
	Extract(c gin.Context) (code.ResultCode, error)
	Exec(ctx context.Context) interface{}
}

var ctxConvert func(c *gin.Context) (newCtx context.Context, err error) = func(c *gin.Context) (newCtx context.Context, err error) {
	return context.Background(), nil
}

// Init 将gin.Context转换成context.Context
func Init(convert func(c *gin.Context) (newCtx context.Context, err error)) {
	ctxConvert = convert
}

// Handler .
func Handler(process Process) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := process.New()
		ctx, err := ctxGet(c)
		if err != nil {
			c.JSON(http.StatusOK, response.NewResponse(ctx, code.Default, err))
			return
		}
		resultCode, err = req.Extract(c)
		if err != nil {
			c.JSON(http.StatusOK, response.NewResponse(ctx, resultCode, err))
			return
		}
		data = req.Exec(ctx)
		c.JSON(http.StatusOK, data)
	}
}
