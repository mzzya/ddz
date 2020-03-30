package ginext

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/util"
)

// Config 配置
type Config interface {
	//ResultInfo code对应info信息
	//todo i18n
	ResultInfo(util.ResultCode) string
	//ContextConvert gin.Context转换为ctx方法，主要作用 可解析用户信息并向下透传给日志组件等记录用户信息
	ContextConvert(c *gin.Context) (newCtx context.Context, err error)
	//ExtractAuthorization 提取授权信息 jwt,token等
	ExtractAuthorization(c *gin.Context) (authorizationInfo interface{}, code util.ResultCode, err error)
}

var conf Config

// Init 初始化
func Init(c Config) {
	conf = c
}

// NullConfig .
type NullConfig struct {
}

// ResultInfo .
func (n *NullConfig) ResultInfo(util.ResultCode) string {
	return ""
}

// ContextConvert .
func (n *NullConfig) ContextConvert(c *gin.Context) (newCtx context.Context, err error) {
	return context.Background(), nil
}

func (n *NullConfig) ExtractAuthorization(c *gin.Context) (authorizationInfo interface{}, code util.ResultCode, err error) {
	return
}
