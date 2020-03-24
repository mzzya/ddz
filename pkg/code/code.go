package code

// ResultCode 结果码
type ResultCode string

const (
	// Success 成功
	Success ResultCode = "success"
	// Default 默认
	Default ResultCode = ""
	//ErrorParameterIsIncorrect 参数不正确
	ErrorParameterIsIncorrect ResultCode = "ErrorParameterIsIncorrect"
)
