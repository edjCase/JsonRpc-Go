package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime/debug"

	"github.com/Gekctek/JsonRpc-Go/pkg/common"
)

func BuildHttpHandler(handler RequestHandler) http.HandlerFunc {
	if handler == nil {
		panic("No json-rpc request handler provided")
	}
	var logger = log.New(os.Stdout, "", 0)
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				logError(logger, r)
				w.WriteHeader(http.StatusInternalServerError)
			}
		}()
		//TODO should it be closed or not?
		defer r.Body.Close()

		bytes := HandleRequest(r.Body, handler, logger)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, writeErr := w.Write(bytes)
		if writeErr != nil {
			panic(common.NestedError(writeErr, "Unable to write the json-rpc response"))
		}
	}
}
func HandleRequest(r io.Reader, handler RequestHandler, logger *log.Logger) []byte {
	bytes, handleErr := handleRequest(r, handler)
	if handleErr == nil {
		return bytes
	}
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
	return []byte(fmt.Sprintf("{\"jsonrpc\":\"2.0\",\"result\":null,\"error\":{\"code\":%d,\"message\":\"%s\",\"data\":%s}}", code, message, data))
}

func handleRequest(r io.Reader, handler RequestHandler) (bytes []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				err = errors.New("Unknown error")
			}
		}
	}()
	info := ParseReader(r)
	response := handler(info)
	response.JsonRpc = "2.0"
	bytes, serializationErr := json.Marshal(response)
	if serializationErr != nil {
		return nil, createError(InternalError, "Unable to serialize the json-rpc response", serializationErr)
	}
	return bytes, nil
}

func logMessage(l *log.Logger, message string) {
	l.Print(message)
}
func logError(l *log.Logger, err interface{}) {
	l.Printf("Error: %s\nStackTrace: %s", err, string(debug.Stack()))
}
