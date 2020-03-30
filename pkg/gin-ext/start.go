package ginext

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/logger"
	"github.com/hellojqk/simple/pkg/util"
	"go.uber.org/zap"
)

var startTime = util.TimeFormat(time.Now())

var hostName = os.Getenv("HOSTNAME")

// getURL 获取url地址
func getURL(addr string) string {
	if strings.Index(addr, ":") == 0 {
		return "http://localhost" + addr
	}

	if strings.Index(addr, "http") != 0 {
		return "http://" + addr
	}
	return addr
}

// Start web服务启动方法 addr 地址 closeWaitTime 关闭服务等待超时限制
func Start(addr string, closeWaitTime time.Duration, engine *gin.Engine) {
	engine.GET("/v", func(c *gin.Context) {
		c.String(http.StatusOK, "进程启动时间:%s\nHOSTNAME:%s", startTime, hostName)
	})
	//初始化gin
	srv := &http.Server{
		Addr:    addr,
		Handler: engine,
	}
	quit := make(chan os.Signal, 1)
	go func() {
		fmt.Println("地址:", getURL(addr))
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			//发生异常主动关闭
			quit <- os.Kill
		}
	}()

	startTime := time.Now()
	logger.Logger.Info("服务已启动", zap.String("server_start_time", util.TimeFormat(startTime)))
	//接收关闭通知 ctrl+c 或 kill
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	stopBeginTime := time.Now()
	//开始关闭服务
	logger.Logger.Info("服务关闭中", zap.String("server_stop_begin_time", util.TimeFormat(stopBeginTime)))

	ctx, cancel := context.WithTimeout(context.Background(), closeWaitTime)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Logger.Error("服务关闭异常", zap.Error(err),
			zap.String("server_start_time", util.TimeFormat(startTime)),
			zap.String("server_stop_begin_time", util.TimeFormat(stopBeginTime)),
			zap.String("server_stop_error_time", util.TimeFormat(time.Now())),
			zap.String("server_live_time", time.Since(startTime).String()),
		)

	}
	log.Println(time.Since(startTime))
	logger.Logger.Info("服务已关闭",
		zap.String("server_start_time", util.TimeFormat(startTime)),
		zap.String("server_stop_begin_time", util.TimeFormat(stopBeginTime)),
		zap.String("server_stop_end_time", util.TimeFormat(time.Now())),
		zap.String("server_live_time", time.Since(startTime).String()),
	)
}

// MergeStart 将多个gin.engine合并至一个gin.engine中启动 默认500个路由，路由数为预估数量 小于实际数量会报错
func MergeStart(addr string, closeWaitTime time.Duration, engines ...*gin.Engine) {
	MergeStartWithRouterCount(addr, closeWaitTime, 500, engines...)
}

// MergeStartWithRouterCount 指定路由数量启动，小于实际数量会报错
func MergeStartWithRouterCount(addr string, closeWaitTime time.Duration, routersLen int, engines ...*gin.Engine) {

	var routerCount int
	routers := make([]gin.RouteInfo, routersLen)
	for _, engine := range engines {
		copy(routers[routerCount:], []gin.RouteInfo(engine.Routes()))
		routerCount += len(engine.Routes())
	}
	engine := gin.New()
	for _, router := range routers {
		if router.Method == "" {
			break
		}
		fmt.Printf("%s\t%s\t%s\n", router.Method, router.Path, router.Handler)
		engine.Handle(router.Method, router.Path, router.HandlerFunc)
	}
	Start(addr, closeWaitTime, engine)
}
