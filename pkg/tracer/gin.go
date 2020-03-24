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

// InjectSpanToGinContext .
func InjectSpanToGinContext(c *gin.Context, span opentracing.Span) {
	if !Enable {
		return
	}
	c.Set(TraceGinSpanKey, span)
	return
}

// ExtractSpanFromGinContext .
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

// GinContextConvert .
func GinContextConvert(ctx context.Context, c *gin.Context) (newCtx context.Context, err error) {
	span, err := ExtractSpanFromGinContext(c)
	if err != nil {
		return ctx, err
	}
	newCtx = opentracing.ContextWithSpan(ctx, span)
	return newCtx, nil
}
