package router

import (
	"fmt"
	"runtime/debug"
)

const (
	ParserError    = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

type Error struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %s", e.Code, e.Message, e.Data)
}

func createError(code int, message string, err error) error {
	return createErrorWithData(code, message, err, nil)
}

func createErrorWithData(code int, message string, err error, data interface{}) error {
	var m string
	if err != nil {
		m = fmt.Sprintf("%s, Inner error: %s, StackTrace: %s", message, err.Error(), string(debug.Stack()))
	} else {
		m = message
	}
	return &Error{Code: code, Message: m, Data: data}
}
