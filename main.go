package main

import (
	"fmt"
	"net/http"

	"github.com/jakubruminski/FYP/go/api/query"
	"github.com/jakubruminski/FYP/go/router/mux"

	"github.com/jakubruminski/FYP/go/utils/logger"
)


func main() {
	fmt.Println("\033[H\033[2J")
	logger := &logger.Logger{}

	port, mux, ok := mux.INIT(logger)
	if !ok { fmt.Println("Failed to initialize router"); return }
	
    ok = query.INITIALISE_DATABASE(logger)
	if !ok { fmt.Println("Failed to initialize database"); return }
	
	fmt.Println("Listening on port http://localhost:" + port)

	http.ListenAndServe(":"+port, mux)
}