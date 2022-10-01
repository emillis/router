package veryFastRouter

import (
	"net/http"
)

//===========[STRUCTS]====================================================================================================

//httpStatusCodeHandlers defines structure of http status code handlers, e.g. 404, 405, etc..
type httpStatusCodeHandlers struct {
	//allowedHttpCodes defines default handlers for various http status codes
	handlers map[int]HandlerFunc
}

//===========[FUNCTIONALITY]====================================================================================================

//newCustomHttpCodeHandlers initializes and returns pointer to a new httpStatusCodeHandlers
func newCustomHttpCodeHandlers() httpStatusCodeHandlers {
	return httpStatusCodeHandlers{
		handlers: map[int]HandlerFunc{
			http.StatusNotFound: func(w http.ResponseWriter, r *http.Request, ai *AdditionalInfo) {
				w.WriteHeader(404)
			},
			http.StatusMethodNotAllowed: func(w http.ResponseWriter, r *http.Request, ai *AdditionalInfo) {
				w.WriteHeader(405)
			},
		},
	}
}
