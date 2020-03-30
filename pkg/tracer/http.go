package tracer

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

//DoBefore http 请求之前
func DoBefore(ctx context.Context, req *http.Request, operationName string) (span opentracing.Span, newCtx context.Context, err error) {
	if !Enable {
		return
	}
	if operationName == "" {
		operationName = req.URL.Host + req.URL.Path
	}
	span, newCtx = opentracing.StartSpanFromContext(ctx, operationName)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	err = InjectSpanToHeader(req, span)
	if err != nil {
		return span, newCtx, err
	}
	return span, newCtx, nil
}

// DoAfter http 请求之后
func DoAfter(span opentracing.Span, resp *http.Response, err error) {
	if !Enable {
		return
	}
	ext.HTTPStatusutil.Set(span, uint16(resp.StatusCode))
	if err != nil {
		span.LogFields(log.Error(err))
	}
	span.Finish()
}

// ExtractSpanFromHeader 从http请求头中获取链路跟踪信息 主要是给自己暴露的API请求过滤器中使用
func ExtractSpanFromHeader(req *http.Request, operationName string) (span opentracing.Span) {
	if !Enable {
		return
	}
	reqSpanCtx, _ := opentracing.GlobalTracer().Extract(opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if reqSpanCtx != nil {
		// 如果调用系统有则继承
		span = opentracing.StartSpan(operationName, opentracing.ChildOf(reqSpanCtx))
	} else {
		// 如果调用放没有则新建
		span = opentracing.StartSpan(operationName)
	}
	return
}

// InjectSpanToHeader 将Span信息注入到Header中
func InjectSpanToHeader(req *http.Request, span opentracing.Span) (err error) {
	if !Enable {
		return
	}
	return opentracing.GlobalTracer().Inject(span.Context(), opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
}
