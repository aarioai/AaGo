package ae

import (
	"fmt"
	"net/http"
)

// Error 定义标准错误结构
type Error struct {
	Code    int            `json:"code"`
	Msg     string         `json:"msg"`
	Caller  string         `json:"caller,omitempty"`
	Details map[string]any `json:"details,omitempty"`
}

// New 使用错误码和消息创建 Error
func New(code int, msg string) *Error {
	return &Error{
		Code:   code,
		Msg:    msg,
		Caller: Caller(2),
	}
}

// NewCode 使用错误码创建 Error
func NewCode(code int) *Error {
	return &Error{
		Code:   code,
		Msg:    http.StatusText(code),
		Caller: Caller(1),
	}
}

// NewMsg 使用消息创建 Error
func NewMsg(format string, args ...any) *Error {
	return &Error{
		Code:   500,
		Msg:    fmt.Sprintf(format, args...),
		Caller: Caller(1),
	}
}

// NewError 从标准 error 创建 Error
func NewError(err error) *Error {
	if err == nil {
		return nil
	}
	return NewMsg(err.Error()).WithCaller(2)
}

// WithCaller 添加调用者信息
func (e *Error) WithCaller(skip int) *Error {
	if e == nil {
		return nil
	}
	e.Caller = Caller(skip + 1)
	return e
}

// WithDetail 添加详细信息
func (e *Error) WithDetail(key string, value any) *Error {
	if e == nil {
		return nil
	}
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// Text 输出错误信息，最好不要使用 Error，避免跟 error 一致，导致人写的时候发生失误
func (e *Error) Text() string {
	if e == nil {
		return "<nil>"
	}

	if e.Caller != "" {
		return fmt.Sprintf("[%d] %s at %s", e.Code, e.Msg, e.Caller)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Msg)
}

// 状态检查方法
func (e *Error) IsNotFound() bool {
	return e != nil && (e.Code == CodeNotFound || e.Code == CodeNoRows || e.Code == CodeGone)
}

func (e *Error) IsServerError() bool {
	return e != nil && e.Code >= 500 && e.Code <= 599
}

func (e *Error) IsRetryWith() bool {
	return e != nil && e.Code == 449 && e.Msg != ""
}
