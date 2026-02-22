package app_error

import (
	"errors"
	"fmt"

	"go.uber.org/zap"
)

type ErrCode int
type ErrType int

const (
	ErrTypeInput ErrType = iota
	ErrTypeInternal
)

type AppError interface {
	error
	Code() ErrCode
	Msg() string
	Type() ErrType
	ErrorDetail() string
	ErrorField() []zap.Field
	WithError(err error) AppError
}

type InputError struct {
	msg  string  // 详情
	code ErrCode // 业务错误码 2xx 参数错误
	err  error
}

func NewInputError(msg string, code ErrCode, err error) *InputError {
	return &InputError{msg: msg, code: code, err: err}
}

func (a *InputError) Error() string {
	return fmt.Sprintf("[ %d ] [invalid input]: %s", a.code, a.msg)
}

func (a *InputError) ErrorDetail() string {
	return fmt.Sprintf("[ %d ] [invalid input]: %s (root: %+v)", a.code, a.msg, a.err)
}

func (a *InputError) Code() ErrCode {
	return a.code
}

func (a *InputError) Msg() string {
	return a.msg
}

func (a *InputError) Type() ErrType {
	return ErrTypeInput
}

func (a *InputError) ErrorField() []zap.Field {
	return []zap.Field{zap.String("error_type", "input error"), zap.String("error_detail", a.ErrorDetail())}
}

func (a *InputError) Unwrap() error {
	return a.err
}

func (a *InputError) WithError(err error) AppError {
	a.err = err
	return a
}

type InternalError struct {
	code ErrCode // 业务错误码 1xx 内部错误
	err  error
}

func NewInternalError(code ErrCode, err error) *InternalError {
	var internal *InternalError
	if errors.As(err, &internal) {
		return internal
	}
	return &InternalError{code: code, err: err}
}

func (a *InternalError) Error() string {
	return fmt.Sprintf("[ %d ] [internal error]: %s", a.code, a.err.Error())

}

func (a *InternalError) Unwrap() error {
	return a.err
}

func (a *InternalError) Code() ErrCode {
	return a.code
}

func (a *InternalError) Msg() string {
	return "internal error"
}

func (a *InternalError) Type() ErrType {
	return ErrTypeInternal
}

func (a *InternalError) ErrorDetail() string {
	return a.Error()
}

func (a *InternalError) ErrorField() []zap.Field {
	return []zap.Field{zap.String("error_type", "internal error"), zap.String("error_detail", a.ErrorDetail())}
}

func (a *InternalError) WithError(err error) AppError {
	a.err = err
	return a
}
