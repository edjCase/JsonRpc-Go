package router

import (
	"encoding/json"
	"errors"
	"fmt"
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
		defer catch(w, logger)
		info, parseErr := parseRequest(r)
		if parseErr != nil {
			panic(common.NestedError(parseErr, "Unable to parse the request into a json-rpc format"))
		}
		response := handler(info)
		bytes, serializationErr := json.Marshal(response)
		if serializationErr != nil {
			panic(common.NestedError(serializationErr, "Unable to serialize the json-rpc response"))
		}

		w.Header().Set("Content-Type", "application/json")
		_, writeErr := w.Write(bytes)
		if writeErr != nil {
			panic(common.NestedError(serializationErr, "Unable to write the json-rpc response"))
		}
	}
}

func catch(w http.ResponseWriter, logger *log.Logger) {
	if r := recover(); r != nil {
		logger.Printf("Error: %s\nStackTrace: %s", r, string(debug.Stack()))
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//TODO better errors and return a json-rpc message
func parseRequest(r *http.Request) (Request, error) {
	bytes, readErr := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if readErr != nil {
		panic(common.NestedError(readErr, "Unable to read the http request body"))
	}
	if len(bytes) == 0 {
		return Request{}, errors.New("No request body supplied")
	}
	var request Request
	serializationErr := json.Unmarshal(bytes, &request)
	if serializationErr != nil {
		return Request{}, common.NestedError(serializationErr, "Unable to deserialize the http request body")
	}
	return request, nil
}

type Response struct {
	Id     string
	Result interface{}
	Error  Error
}

type Request struct {
	Id         string
	Method     string
	Parameters []interface{}
}

type Error struct {
	Code    int
	Message string
	Data    interface{}
}

type RequestHandler func(Request) Response
