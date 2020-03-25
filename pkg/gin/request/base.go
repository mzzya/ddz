package request

import (
	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/code"
)

// Base 统一返回对象
type Base struct {
}

//Extract 参数提取
func (b *Base) Extract(c *gin.Context) (code.ResultCode, error) {
	return b.DefaultExtract(b, c)
}

//DefaultExtract 默认提取方法
func (b *Base) DefaultExtract(data interface{}, c *gin.Context) (co code.ResultCode, err error) {
	return b.ExtractWithBindFunc(data, c.ShouldBind)
}

// ExtractWithBindFunc extract with bindFunc
func (b *Base) ExtractWithBindFunc(data interface{}, bindFunc func(obj interface{}) error) (c code.ResultCode, err error) {
	err = bindFunc(data)
	return
}
