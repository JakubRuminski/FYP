package mux

import (
	"net/http"

	"github.com/jakubruminski/FYP/go/router/request"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func INIT() (port string, mux *http.ServeMux, ok bool) {
	mux = http.NewServeMux()

	logger := logger.Logger{}
	environment, ok := env.Get(&logger, "ENVIRONMENT")
	logger.SetEnvironment(environment)

	mux.HandleFunc("/", logger.Middleware(request.HandleRequest))
	mux.HandleFunc("/static/", logger.Middleware(request.HandleRequest))

	mux.HandleFunc("/api/search", logger.Middleware(request.HandleApiRequest))

	port, ok = env.Get(&logger, "PORT")
	if !ok {
		return "", nil, false
	}

	return port, mux, true
}