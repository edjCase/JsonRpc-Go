package main

import (
	"github.com/Gekctek/JsonRpc-Go/pkg/router"
)

func main() {
	router.Run("/", 8000, handle)
}

func handle(info router.Request) router.Response {

	return router.Response{Id: "1", Result: 1}
}
