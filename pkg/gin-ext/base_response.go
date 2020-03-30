package ginext

import (
	"context"

	"github.com/hellojqk/simple/pkg/util"
)

// BaseResponse 统一返回对象
type BaseResponse struct {
	Code      util.ResultCode `json:"code"`             //业务状态码
	Status    bool            `json:"status"`           //状态【code==200】
	Msg       string          `json:"msg"`              //提示消息 一般为状态码对应信息
	RequestID interface{}     `json:"req_id,omitempty"` //请求ID
	Desc      string          `json:"desc,omitempty"`   //Error信息
}

// NewResponse 根据业务状态码和err信息创建新的结构返回
func NewResponse(ctx context.Context, resultCode util.ResultCode, err error) BaseResponse {
	if err == nil && (resultCode == util.Success || resultCode == util.Default) {
		return NewSuccessResponse(ctx)
	}
	return BaseResponse{RequestID: 0, Code: resultCode, Status: resultCode == util.Success, Msg: conf.ResultInfo(resultCode), Desc: err.Error()}
}

// NewSuccessResponse success response
func NewSuccessResponse(ctx context.Context) BaseResponse {
	return BaseResponse{RequestID: 0, Code: util.Success, Status: true, Msg: conf.ResultInfo(util.Success)}
}
