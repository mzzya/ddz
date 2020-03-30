package router

import (
	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/inertnal/controller"
	ginext "github.com/hellojqk/simple/pkg/gin-ext"
)

// NewRouter 初始化gin引擎，路由分组，参数验证插件V9版本
func NewRouter() *gin.Engine {
	ginext.Init(&ginext.NullConfig{})

	engine := gin.New()
	engine.Use(gin.Recovery())
	//不需要权限验证的放这里
	v1group := engine.Group("/v1")
	{
		v1group.GET("/test", ginext.Handler(&controller.TestRequest{}))
	}
	return engine
}
