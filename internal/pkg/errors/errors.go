package errors

import "fmt"

// ErrorCode 错误码类型
type ErrorCode int

const (
	// 系统级错误码 (1000-1999)
	ErrInvalidParam  ErrorCode = 1000
	ErrUnauthorized  ErrorCode = 1001
	ErrInvalidSign   ErrorCode = 1002
	ErrExpired       ErrorCode = 1003
	ErrInternal      ErrorCode = 1004
	ErrTooManyReq    ErrorCode = 1005
	ErrServiceUnavailable ErrorCode = 1006

	// 设备相关错误码 (2000-2999)
	ErrDeviceNotFound    ErrorCode = 2000
	ErrDeviceOffline     ErrorCode = 2001
	ErrDeviceTypeInvalid ErrorCode = 2002
	ErrInvalidSN         ErrorCode = 2003
	ErrCacheFull         ErrorCode = 2004
	ErrFirmwareVersion   ErrorCode = 2005
)

// Error 自定义错误类型
type Error struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("error: code=%d, message=%s", e.Code, e.Message)
}

// New 创建新的错误
func New(code ErrorCode, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
} 