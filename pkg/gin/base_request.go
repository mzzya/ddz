package gin

import (
	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/code"
)

// BaseRequest 统一返回对象
type BaseRequest struct {
}

//Extract 参数提取
func (b *BaseRequest) Extract(c *gin.Context) (code.ResultCode, error) {
	return b.DefaultExtract(b, c)
}

//DefaultExtract 默认提取方法
func (b *BaseRequest) DefaultExtract(data interface{}, c *gin.Context) (co code.ResultCode, err error) {
	return b.ExtractWithBindFunc(data, c.ShouldBind)
}

// ExtractWithBindFunc extract with bindFunc
func (b *BaseRequest) ExtractWithBindFunc(data interface{}, bindFunc func(obj interface{}) error) (c code.ResultCode, err error) {
	err = bindFunc(data)
	return
}
