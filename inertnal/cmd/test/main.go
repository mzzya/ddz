package main

import (
	"time"

	"github.com/hellojqk/simple/inertnal/router"
	ginext "github.com/hellojqk/simple/pkg/gin-ext"
)

func main() {
	ginext.MergeStart(":8080", time.Minute, router.NewRouter())
}
