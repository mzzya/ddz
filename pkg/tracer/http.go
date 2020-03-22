package tracer

import (
	"context"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
)

//DoBefore http 请求之前
func DoBefore(ctx context.Context, req *http.Request, operationName string) (span opentracing.Span, err error) {
	if operationName == "" {
		operationName = req.URL.Host + req.URL.Path
	}
	span, _ = opentracing.StartSpanFromContext(ctx, operationName)
	ext.HTTPMethod.Set(span, req.Method)
	ext.HTTPUrl.Set(span, req.URL.String())
	err = span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header))
	if err != nil {
		return nil, err
	}
	return span, nil
}

// DoAfter http 请求之后
func DoAfter(span opentracing.Span, resp *http.Response, err error) {
	ext.HTTPStatusCode.Set(span, uint16(resp.StatusCode))
	if err != nil {
		span.LogFields(log.Error(err))
	}
	span.Finish()
}
