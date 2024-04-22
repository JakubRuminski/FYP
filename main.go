package main

import (
	"fmt"
	"net/http"

	"github.com/jakubruminski/FYP/go/api/query"
	"github.com/jakubruminski/FYP/go/router/mux"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func main() {
	fmt.Println("\033[H\033[2J")
	logger := &logger.Logger{}

	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { logger.ERROR("Failed to get DB_AVAILABLE from environment"); return }

	if db_available {
    	ok := query.INITIALISE_DATABASE(logger)
		if !ok { logger.ERROR("Failed to initialize database"); return }
	}
	
	port, mux, ok := mux.INIT(logger)
	if !ok { logger.ERROR("Failed to initialize router"); return }

	logger.INFO("Listening on port http://localhost:" + port)

	http.ListenAndServe(":"+port, mux)
}