package router

import (
	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/inertnal/controller"
	ginext "github.com/hellojqk/simple/pkg/gin-ext"
)

// NewRouter 初始化gin引擎，路由分组，参数验证插件V9版本
func NewRouter() *gin.Engine {
	ginext.Init(&ginext.NullConfig{})
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gzip.Gzip(gzip.DefaultCompression))
	engine.Use(ginext.CORS(), ginext.Recovery(), ginext.Logger())
	//不需要权限验证的放这里
	v1group := engine.Group("/v1")
	{
		v1group.GET("/test", ginext.Handler(&controller.TestRequest{}))
	}
	return engine
}
