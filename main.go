package main

import (
	"fmt"
	"time"

	"github.com/hellojqk/simple/inertnal/router"
	"github.com/hellojqk/simple/pkg/config"
	ginext "github.com/hellojqk/simple/pkg/gin-ext"
	"github.com/hellojqk/simple/pkg/logger"
	"github.com/hellojqk/simple/pkg/tracer"
	"github.com/hellojqk/simple/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

// 初始化
func init() {
	DI()
}

// DI 依赖注入
func DI() {
	c := dig.New()
	var err error
	if err = c.Provide(config.DefaultViper); err != nil {
		fmt.Print(errors.WithMessage(err, "Provide DefaultViper"))
	}
	err = c.Invoke(logger.Init)
	if err != nil {
		fmt.Print(errors.WithMessage(err, "Invoke"))
	}
	err = c.Invoke(tracer.Init)
	if err != nil {
		fmt.Print(errors.WithMessage(err, "Invoke"))
	}
}

// Close 关闭
func Close() {
	util.Close()
}

func main() {
	logger.Logger.Info("app start")
	ginext.MergeStart(":8080", time.Minute, router.NewRouter())
	Close()
	logger.Logger.Info("app stop")
}
