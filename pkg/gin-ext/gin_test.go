package ginext

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/gin-gonic/gin"
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

func TestMain(m *testing.M) {
	Init(&NullConfig{})
	router = gin.New()
	router.Use(gin.Recovery())
	router.GET("/", Handler(&TestRequest{}))
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
}

func TestGin(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.NotEqual(t, "", w.Body.String())
	t.Logf("%s\n", w.Body.String())
}

//go test -v -bench . -run ^Gin$  -benchmem -cpuprofile cpu.out -memprofile mem.out
func BenchmarkGin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/", nil)
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
