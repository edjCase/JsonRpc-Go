package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Gekctek/JsonRpc-Go/pkg/common"
)

//Run accepts all requests for the pattern
func Run(pattern string, port int, handler RequestHandler) {
	http.Handle(pattern, baseHandler(handler))
	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func baseHandler(handler RequestHandler) http.HandlerFunc {
	var logger = log.New(os.Stdout, "", 0)
	return func(w http.ResponseWriter, r *http.Request) {
		defer logResponse(w, logger)
		bytes, handleErr := handleRequest(r.Body, handler)
		if handleErr != nil {
			var code int
			var message string
			var data = "null"
			if e, ok := handleErr.(*Error); ok {
				code = e.Code
				message = e.Message
				if e.Data != nil {
					b, jsonErr := json.Marshal(e.Data)
					logError(logger, jsonErr)
					if b != nil {
						data = fmt.Sprintf("\"%s\"", string(b))
					}
				}
			} else {
				code = InternalError
				message = "An unknown error has occurred"
			}
			bytes = []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"result\":null,\"error\":{\"code\":%d,\"message\":\"%s\",\"data\":%s}}", code, message, data))
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, writeErr := w.Write(bytes)
		if writeErr != nil {
			panic(common.NestedError(writeErr, "Unable to write the json-rpc response"))
		}
	}
}

func handleRequest(r io.ReadCloser, handler RequestHandler) (bytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown panic")
			}
		}
	}()
	info := parseRequest(r)
	response := handler(info)
	response.JsonRpc = "2.0"
	bytes, serializationErr := json.Marshal(response)
	if serializationErr != nil {
		return nil, createError(InternalError, "Unable to serialize the json-rpc response", serializationErr)
	}
	return bytes, nil
}

func logResponse(w http.ResponseWriter, logger *log.Logger) {
	if r := recover(); r != nil {
		logError(logger, r)
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func parseRequest(r io.ReadCloser) Request {
	bytes, readErr := ioutil.ReadAll(r)
	defer r.Close()
	if readErr != nil {
		panic(createError(InternalError, "Unable to read the request body", readErr))
	}
	if len(bytes) == 0 {
		panic(createError(InvalidRequest, "No request body supplied", nil))
	}
	var request Request
	serializationErr := json.Unmarshal(bytes, &request)
	if serializationErr != nil {
		panic(createError(ParserError, "Unable to deserialize request", serializationErr))
	}
	return request
}

const (
	ParserError    = -32700
	InvalidRequest = -32600
	MethodNotFound = -32601
	InvalidParams  = -32602
	InternalError  = -32603
)

func logMessage(l *log.Logger, message string) {
	l.Print(message)
}
func logError(l *log.Logger, err interface{}) {
	l.Printf("Error: %s\nStackTrace: %s", err, string(debug.Stack()))
}

type Response struct {
	JsonRpc string      `json:"jsonrpc"`
	Id      string      `json:"id"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error"`
}

type Request struct {
	JsonRpc string
	Id      string
	Method  string
	Params  interface{}
}

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

type RequestHandler func(Request) Response
