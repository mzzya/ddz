package ginext

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/util"
)

// Process 请求封装
type Process interface {
	New() Process
	Extract(c *gin.Context) (util.ResultCode, error)
	Exec(ctx context.Context) interface{}
}

// Handler .
func Handler(process Process) gin.HandlerFunc {
	return func(c *gin.Context) {
		req := process.New()
		ctx, err := conf.ContextConvert(c)
		if err != nil {
			c.JSON(http.StatusOK, NewResponse(ctx, util.Default, err))
			return
		}
		resultCode, err := req.Extract(c)
		if err != nil {
			c.JSON(http.StatusOK, NewResponse(ctx, resultCode, err))
			return
		}
		data := req.Exec(ctx)
		c.JSON(http.StatusOK, data)
	}
}
