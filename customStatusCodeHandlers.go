package veryFastRouter

import (
	"net/http"
)

//===========[STATIC]====================================================================================================

//Http status codes defined here can be handled with custom handlers by the client.
//Status codes not defined here will not be allowed to have custom handlers.
var allowedHttpStatusCodes = map[int]string{
	401: "401 Unauthorized",
	403: "403 Forbidden",
	404: "404 Not Found",
	405: "405 Method Not Allowed",
}

//===========[STRUCTS]====================================================================================================

//httpStatusCodeHandlers defines structure of http status code handlers, e.g. 404, 405, etc..
type httpStatusCodeHandlers struct {
	//allowedHttpCodes defines default handlers for various http status codes
	handlers map[int]HandlerFunc
}

//===========[FUNCTIONALITY]====================================================================================================

//newCustomHttpCodeHandlers initializes and returns pointer to a new httpStatusCodeHandlers
func newCustomHttpCodeHandlers() httpStatusCodeHandlers {
	h := httpStatusCodeHandlers{
		handlers: make(map[int]HandlerFunc, len(allowedHttpStatusCodes)),
	}

	for k, v := range allowedHttpStatusCodes {
		h.handlers[k] = func(w http.ResponseWriter, r *http.Request, ai *AdditionalInfo) {
			w.WriteHeader(k)
			w.Write([]byte(v))
		}
	}

	return h
}
