package gin

import "github.com/gin-gonic/gin"

// DefaultGin 默认初始化方法
func DefaultGin() *gin.Engine {
	return GinWithHandler()
}

// GinWithHandler 在Use自定义异常和日志中间件前执行的路由绑定
func GinWithHandler(handler ...HandlerFunc) *gin.Engine {
	return nil
}
