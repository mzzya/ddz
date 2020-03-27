package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/code"
)

// Config 配置
type Config interface {
	ResultInfo(code.ResultCode) string
	ContextConvert(c *gin.Context) (newCtx context.Context, err error)
}

var conf Config

// Init 初始化
func Init(c Config) {
	conf = c
}

// NullConfig .
type NullConfig struct{

}
// ResultInfo .
func (n *NullConfig)ResultInfo(code.ResultCode) string{
	return ""
}

// ContextConvert .
func (n *NullConfig)ContextConvert(c *gin.Context) (newCtx context.Context, err error){
	return context.Background(),nil
}