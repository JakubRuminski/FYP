package searchs

import (
	
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "searches"

type SearchTerms struct {
	SearchTerm          string     `json:"search_term"`
	ProductID		    int        `json:"product_id"`

	FetchCount          int        `json:"fetch_count"`
	LastFetch           int        `json:"last_fetch"`
	ProductExists       bool       `json:"product_exists"`
}

func INIT(logger *logger.Logger) (ok bool) {
	
	query := `
		CREATE TABLE IF NOT EXISTS searches
		(
			search_term      VARCHAR(255),
			product_id       INT,
			fetch_count      INT,
			last_fetch       INT,
			product_exists   BOOLEAN
		)
	`

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the searches table")
		return false
	}

	return true	
}