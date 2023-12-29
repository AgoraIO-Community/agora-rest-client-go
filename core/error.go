package core

import (
	"fmt"
)

type InternalErr struct {
	Err string
}

func (i *InternalErr) Error() string {
	return i.Err
}

func NewInternalErr(msg string) *InternalErr {
	return &InternalErr{Err: msg}
}

type RetryErr struct {
	needRetry bool
	Err       error
}

func (r *RetryErr) Error() string {
	return r.Err.Error()
}

func (r *RetryErr) Unwrap() error {
	return r.Err
}

func NewRetryErr(needRetry bool, err error) *RetryErr {
	return &RetryErr{
		needRetry: needRetry,
		Err:       err,
	}
}

func (r *RetryErr) NeedRetry() bool {
	return r.needRetry
}

type GatewayErr struct {
	Code int
	Msg  string
}

func NewGatewayErr(code int, msg string) *GatewayErr {
	return &GatewayErr{Code: code, Msg: msg}
}

func (g *GatewayErr) Error() string {
	return fmt.Sprintf("statusCode:%d,body:%s", g.Code, g.Msg)
}
