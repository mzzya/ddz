package ginext

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/hellojqk/simple/pkg/logger"
	"github.com/hellojqk/simple/pkg/tracer"
	"github.com/hellojqk/simple/pkg/util"
	"github.com/opentracing/opentracing-go/ext"
	tracerLog "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// CORS 跨域配置
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取请求发起域名
		var origin = c.GetHeader("Origin")
		// todo 允许跨域域名 仅支持单个配置 且不支持正则、通配符
		// var env = util.GetAppEnv()
		// switch env {

		// }
		c.Header("Access-Control-Allow-Origin", origin)
		// 允许cookie
		c.Header("Access-Control-Allow-Credentials", "true")
		//跨域请求正式发起前会先发起一个 OPTIONS 请求获取允许的跨域Method和Header
		if c.Request.Method == "OPTIONS" {
			c.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE,PUT")
			c.Header("Access-Control-Allow-Headers", "powercode,token,applicationid,content-type")
			c.Status(200)
			return
		}
	}
}

// Logger 请求日志记录
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		util.PrintJSONWithColor(c.Request.URL)
		fmt.Print("c.Request.RequestURI", c.Request.RequestURI)
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery
		if raw != "" {
			path = path + "?" + raw
		}
		span := tracer.ExtractSpanFromHeader(c.Request, c.Request.URL.Path)
		if span != nil {
			defer span.Finish()
			tracer.InjectSpanToGinContext(c, span)
			span.LogFields(tracerLog.String("query", c.Request.URL.RawQuery))
			ext.HTTPMethod.Set(span, c.Request.Method)
			ext.HTTPUrl.Set(span, path)
		}
		c.Next()
		logger.Logger.Info("request", zap.String("method", c.Request.Method), zap.Int("status", c.Writer.Status()), zap.String("url", path))
		if span != nil {
			ext.HTTPStatusCode.Set(span, uint16(c.Writer.Status()))
		}
	}
}

// Recovery 异常捕获
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Check for a broken connection, as it is not really a
				// condition that warrants a panic stack trace.
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				stack := stack(3)
				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				headers := strings.Split(string(httpRequest), "\r\n")
				for idx, header := range headers {
					current := strings.Split(header, ":")
					if current[0] == "Authorization" {
						headers[idx] = current[0] + ": *"
					}
				}
				if brokenPipe {
					logger.Logger.Error("gin recover brokenpipe", zap.Reflect("error", err), zap.ByteString("request", httpRequest))
				} else {
					logger.Logger.Error("gin recover", zap.Reflect("error", err), zap.ByteString("stack", stack))
				}
				span, spanErr := tracer.ExtractSpanFromGinContext(c)
				if spanErr != nil {
					logger.Logger.Error("gin recover extract span", zap.Reflect("error", err), zap.Error(spanErr), zap.ByteString("stack", stack))
				}
				if span != nil {
					span.SetTag("error", "true")
					span.LogFields(
						tracerLog.Bool("brokenPipe", brokenPipe),
						tracerLog.Object("recover", err))
					span.Finish()
				}

				// If the connection is dead, we can't write a status to it.
				if brokenPipe {
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, BaseResponse{Code: util.Error, Status: false, Msg: conf.ResultInfo(util.Error), Desc: string(stack)})
				}
			}
		}()
		c.Next()
	}
}

// stack returns a nicely formatted stack frame, skipping skip frames.
func stack(skip int) []byte {
	buf := new(bytes.Buffer) // the returned data
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)
		if file != lastFile {
			data, err := ioutil.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}
	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contains dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
