package mux

import (
	"net/http"

	"github.com/jakubruminski/FYP/go/router/request"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


var MAX_REQUESTS int

// This function is a middleware that limits the number of requests
// that can be handled at the same time.
//
// If called with a semaphore of size 1, it will only allow one request to be handled at a time.
// 
func RequestLimiter(logger *logger.Logger, next http.HandlerFunc) http.HandlerFunc {
	semaphore := make(chan struct{}, MAX_REQUESTS)
	logger.DEBUG("Created request limiter")

	return func(w http.ResponseWriter, r *http.Request) {
		semaphore <- struct{}{}
		logger.DEBUG("Acquired semaphore")
		
		defer func() { 
			logger.DEBUG("Released semaphore")
			<-semaphore
		}()

		next(w, r)
		logger.DEBUG("Request handled")
	}
}

func INIT(logger *logger.Logger) (port string, mux *http.ServeMux, ok bool) {
	mux = http.NewServeMux()

	environment := "ENVIRONMENT"
	port = "PORT"
	ok = env.GetKeys(logger, &environment, &port)
	if !ok { return "", nil, false }

	verbose, ok := env.GetBool( logger, "VERBOSE" )
	if !ok { return "", nil, false }

	MAX_REQUESTS, ok = env.GetInt( logger, "MAX_REQUESTS" )
    if !ok { return "", nil, false }
	logger.SetFlags(environment, verbose, "root-logger")

	mux.HandleFunc("/", RequestLimiter( logger, request.HandleRequest ))
	mux.HandleFunc("/basket", RequestLimiter( logger, request.HandleRequest ))
	mux.HandleFunc("/static/", RequestLimiter( logger, request.HandleRequest ))

	mux.HandleFunc("/api/search", RequestLimiter( logger, request.HandleApiRequest ))

	mux.HandleFunc("/api/add_item", RequestLimiter( logger, request.HandleApiRequest ))
	mux.HandleFunc("/api/remove_item", RequestLimiter( logger, request.HandleApiRequest ))

	mux.HandleFunc("/api/get_items", RequestLimiter( logger, request.HandleApiRequest ))

	return port, mux, true
}