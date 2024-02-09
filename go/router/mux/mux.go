package mux

import (
	"net/http"

	"github.com/jakubruminski/FYP/go/router/request"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func INIT(logger *logger.Logger) (port string, mux *http.ServeMux, ok bool) {
	mux = http.NewServeMux()

	environment := "ENVIRONMENT"
	port = "PORT"
	ok = env.GetKeys(logger, &environment, &port)
	if !ok { return "", nil, false }
	
	logger.SetEnvironment(environment)

	mux.HandleFunc("/", logger.Middleware(request.HandleRequest))
	mux.HandleFunc("/static/", logger.Middleware(request.HandleRequest))

	mux.HandleFunc("/api/search", logger.Middleware(request.HandleApiRequest))
	mux.HandleFunc("/api/add_item", logger.Middleware(request.HandleApiRequest))

	return port, mux, true
}