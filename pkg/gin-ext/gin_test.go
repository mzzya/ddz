package ginext

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/config"
	"github.com/hellojqk/simple/pkg/logger"
	"github.com/hellojqk/simple/pkg/tracer"
	"github.com/hellojqk/simple/pkg/util"
	"github.com/pkg/errors"
	"go.uber.org/dig"
	"gopkg.in/go-playground/assert.v1"
)

type TestResponse struct {
	BaseResponse
}

var requestPool = sync.Pool{
	New: func() interface{} {
		return &TestRequest{}
	},
}

type TestRequest struct {
	BaseRequest
}

func (r *TestRequest) New() Process {
	return &TestRequest{}
	// return requestPool.Get().(*TestRequest)
}
func (r *TestRequest) Exec(ctx context.Context) interface{} {
	resp := TestResponse{}
	resp.BaseResponse = NewSuccessResponse(ctx)
	// requestPool.Put(r)
	return resp
}

var router *gin.Engine

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
func TestMain(m *testing.M) {
	DI()
	Init(&NullConfig{})
	router = gin.New()
	gin.SetMode(gin.ReleaseMode)
	router.Use(CORS(), Recovery(), Logger())
	router.GET("/test", Handler(&TestRequest{}))
	router.GET("/1", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/2", func(c *gin.Context) {
		c.ShouldBind(&TestRequest{})
		NewSuccessResponse(context.Background())
		c.Status(http.StatusOK)
	})
	router.GET("/3", func(c *gin.Context) {
		c.JSON(http.StatusOK, BaseResponse{})
	})
	router.GET("/4", func(c *gin.Context) {
		c.ShouldBind(&TestRequest{})
		resp := TestResponse{BaseResponse: NewSuccessResponse(context.Background())}
		c.JSON(http.StatusOK, resp)
	})
	m.Run()
	util.Close()
}

func TestGin(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test?a=b", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
	t.Logf("%s\n", w.Body.String())
}

//go test -v -bench . -run ^Gin$  -benchmem -cpuprofile cpu.out -memprofile mem.out
func BenchmarkGin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/test", nil)
		router.ServeHTTP(w, req)
	}
}

func BenchmarkGin1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/1", nil)
		router.ServeHTTP(w, req)
	}
}
func BenchmarkGin2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/2", nil)
		router.ServeHTTP(w, req)
	}
}
func BenchmarkGin3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/3", nil)
		router.ServeHTTP(w, req)
	}
}
func BenchmarkGin4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/4", nil)
		router.ServeHTTP(w, req)
	}
}
