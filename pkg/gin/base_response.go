package gin

import (
	"context"

	"github.com/hellojqk/simple/pkg/code"
)

// BaseResponse 统一返回对象
type BaseResponse struct {
	Code      code.ResultCode `json:"code"`             //业务状态码
	Status    bool            `json:"status"`           //状态【code==200】
	Msg       string          `json:"msg"`              //提示消息 一般为状态码对应信息
	RequestID interface{}     `json:"req_id,omitempty"` //请求ID
	Desc      string          `json:"desc,omitempty"`   //Error信息
}

// NewResponse 根据业务状态码和err信息创建新的结构返回
func NewResponse(ctx context.Context, resultCode code.ResultCode, err error) BaseResponse {
	if err == nil && (resultCode == code.Success || resultCode == code.Default) {
		return NewSuccessResponse(ctx)
	}
	return BaseResponse{RequestID: 0, Code: resultCode, Status: resultCode == code.Success, Msg: conf.ResultInfo(resultCode), Desc: err.Error()}
}

// NewSuccessResponse success response
func NewSuccessResponse(ctx context.Context) BaseResponse {
	return BaseResponse{RequestID: 0, Code: code.Success, Status: true, Msg: conf.ResultInfo(code.Success)}
}
