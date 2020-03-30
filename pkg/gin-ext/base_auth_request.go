package ginext

import (
	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/util"
)

// BaseAuthRequest 统一返回对象
type BaseAuthRequest struct {
	BaseRequest
	AuthInfo interface{}
}

//Extract 参数提取
func (b *BaseAuthRequest) Extract(c *gin.Context) (util.ResultCode, error) {
	return b.DefaultExtract(b, c)
}
