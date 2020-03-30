package util

// ResultCode 结果码
type ResultCode string

const (
	// Success 成功
	Success ResultCode = "success"
	// Error 异常
	Error ResultCode = "error"
	// Default 默认
	Default ResultCode = ""
	//ErrorParameterIsIncorrect 参数不正确
	ErrorParameterIsIncorrect ResultCode = "ErrorParameterIsIncorrect"
)
