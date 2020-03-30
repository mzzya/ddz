package tracer

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
)

const (
	//TraceGinSpanKey opentracing在gin中保存span的Key
	TraceGinSpanKey = "tracer.span"
)

// GinHandler .
func GinHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		span := ExtractSpanFromHeader(c.Request, "")
		if span != nil {
			defer span.Finish()
		}
		c.Next()
	}
}

// InjectSpanToGinContext 将span写入ginContext
func InjectSpanToGinContext(c *gin.Context, span opentracing.Span) {
	if !Enable || span == nil {
		return
	}
	c.Set(TraceGinSpanKey, span)
	return
}

// ExtractSpanFromGinContext 从ginContext提取span信息
func ExtractSpanFromGinContext(c *gin.Context) (span opentracing.Span, err error) {
	if !Enable {
		return
	}
	iSpan, iExists := c.Get(TraceGinSpanKey)
	if !iExists {
		return nil, err
	}
	span, ok := iSpan.(opentracing.Span)
	if !ok {
		return nil, nil
	}
	return span, nil
}

// GinContextConvert ginContext转换成ctx
func GinContextConvert(ctx context.Context, c *gin.Context) (newCtx context.Context, err error) {
	span, err := ExtractSpanFromGinContext(c)
	if err != nil {
		return ctx, err
	}
	newCtx = opentracing.ContextWithSpan(ctx, span)
	return newCtx, nil
}
