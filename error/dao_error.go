package error

import (
	"fmt"
	"reflect"
	"time"
)

type ConnError struct {
	ConnTime time.Time
	Addr     string
	Password string
}

func (e *ConnError) Error() string {
	return "At " + e.ConnTime.String() + "  Connect Error ,Check your Address or Password " + e.Addr
}

type TypeAssertError[V any, T any] struct {
	ValType T
	Value   V
}

func (e *TypeAssertError[T, V]) Error() string {
	return "TypeAssertion Wrong, Check Types, The Type Aim To be Assert To " + reflect.TypeOf(e.Value).String()
}

func NewTypeAssertError[V any, T any](value V) *TypeAssertError[V, T] {
	return &TypeAssertError[V, T]{Value: value}
}

// CloseError 关闭数据库的错误包装
type CloseError struct {
	Operation string
	Err       error
}

func (e *CloseError) Error() string {
	return fmt.Sprintf("%s operation failed: %v", e.Operation, e.Err)
}
