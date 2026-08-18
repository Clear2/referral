package errors

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError - format error for API response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Err     error  `json:"-"`
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(code int, message string) *APIError {
	return &APIError{Code: code, Message: message}
}

func WrapAPIError(base *APIError, err error) *APIError {
	if err == nil {
		return &APIError{Code: base.Code, Message: base.Message}
	}
	return &APIError{
		Code:    base.Code,
		Message: base.Message,
		Err:     err,
	}
}

func FromError(err error) *APIError {
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		return apiErr
	}
	return &APIError{
		Code:    http.StatusInternalServerError,
		Message: "Internal Server Error",
		Err:     err,
	}
}

var (
	ErrBadRequest          = NewAPIError(http.StatusBadRequest, "请求参数不正确")
	ErrUnauthorized        = NewAPIError(http.StatusUnauthorized, "账号或密码错误")
	ErrSessionUnauthorized = NewAPIError(http.StatusUnauthorized, "未登录或登录状态已失效")
	ErrForbidden           = NewAPIError(http.StatusForbidden, "没有操作权限")
	ErrNotFound            = NewAPIError(http.StatusNotFound, "请求的资源不存在")
	ErrConflict            = NewAPIError(http.StatusConflict, "数据已存在")
	ErrInternalServerError = NewAPIError(http.StatusInternalServerError, "服务器内部错误")
)
