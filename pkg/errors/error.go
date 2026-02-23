package errors

// 业务错误码
const (
	CodeOK            = 200
	CodeBadRequest    = 400
	CodeUnauthorized  = 401
	CodeForbidden     = 403
	CodeNotFound      = 404
	CodeInternalError = 500
)

// AppError 应用错误
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

// New 创建错误
func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// Wrap 包装错误
func Wrap(err error, message string) *AppError {
	return &AppError{
		Code:    CodeInternalError,
		Message: message,
		Err:     err,
	}
}

// 预定义错误
var (
	ErrBadRequest     = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized   = New(CodeUnauthorized, "未授权")
	ErrForbidden      = New(CodeForbidden, "禁止访问")
	ErrNotFound       = New(CodeNotFound, "资源不存在")
	ErrInternalServer = New(CodeInternalError, "服务器内部错误")
)

// 业务错误
var (
	ErrImageUploadFailed  = New(5001, "图片上传失败")
	ErrModelServiceFailed = New(5002, "模型服务调用失败")
	ErrRecordSaveFailed   = New(5003, "记录保存失败")
)
