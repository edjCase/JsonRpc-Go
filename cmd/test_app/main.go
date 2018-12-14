package main

import (
	"net/http"

	"github.com/Gekctek/JsonRpc-Go/pkg/router"
)

func main() {
	http.Handle("/", router.BuildHttpHandler(handle))
	http.ListenAndServe(":8000", nil)
}

func handle(info router.Request) router.Response {
	return router.Response{Id: info.Id, Result: info.Method, Error: &router.Error{Code: -1, Data: info.Params}}
}
