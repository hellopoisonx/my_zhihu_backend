package app_error

import (
	"fmt"
)

type InputError struct {
	Msg  string // 详情
	Code int    // 业务错误码 2xx 参数错误
	Err  error
}

func NewInputError(msg string, code int, err error) *InputError {
	return &InputError{Msg: msg, Code: code, Err: err}
}

func (a *InputError) Error() string {
	category := "invalid input"
	if a.Err != nil {
		return fmt.Sprintf("[ %d ] %s: %s (root: %+#v)", a.Code, category, a.Msg, a.Err)
	}
	return fmt.Sprintf("[ %d ] %s: %s", a.Code, category, a.Msg)
}

func (a *InputError) Unwrap() error {
	return a.Err
}

type InternalError struct {
	Code int // 业务错误码 1xx 内部错误
	Err  error
}

func NewInternalError(code int, err error) *InternalError {
	return &InternalError{Code: code, Err: err}
}

func (a *InternalError) Error() string {
	category := "internal error"
	if a.Err != nil {
		return fmt.Sprintf("[ %d ] %s: %s", a.Code, category, a.Err.Error())
	}
	return fmt.Sprintf("[ %d ] %s: %s", a.Code, category, "unknown")
}

func (a *InternalError) Unwrap() error {
	return a.Err
}
