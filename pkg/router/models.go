package router

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

type RequestHandler func(Request) Response
