package router

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

func ParseReader(r io.Reader) Request {
	bytes, readErr := ioutil.ReadAll(r)
	if readErr != nil {
		panic(createError(InternalError, "Unable to read the request body", readErr))
	}
	return ParseBytes(bytes)
}

func ParseBytes(bytes []byte) Request {
	if len(bytes) == 0 {
		panic(createError(InvalidRequest, "Unable to deserialize request, it is empty", nil))
	}
	var request Request
	serializationErr := json.Unmarshal(bytes, &request)
	if serializationErr != nil {
		panic(createError(ParserError, "Unable to deserialize request", serializationErr))
	}
	return request
}
