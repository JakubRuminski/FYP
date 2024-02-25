package query_searchs

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "searches"

type SearchTerm struct {
	SearchTerm          string     `json:"search_term"`
	ProductID		    int        `json:"product_id"`

	FetchCount          int        `json:"fetch_count"`
	LastFetch           int        `json:"last_fetch"`
}

func INIT(logger *logger.Logger) (ok bool) {
	
	query := `
		CREATE TABLE IF NOT EXISTS searches
		(
			search_term      VARCHAR(255),
			product_id       INT,
			fetch_count      INT,
			last_fetch       INT
		)
	`

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the searches table")
		return false
	}

	return true	
}

func Get(logger *logger.Logger, tx *sql.Tx, searchTerm string) (productIDs *[]*string, ok bool) {
	searchTerm = strings.ToLower(searchTerm)

	query := `SELECT product_id FROM searches WHERE search_term = $1`

	productIDs = &[]*string{}
	ok = postgres.ExecuteContextLookUpQuery(logger, tx, get, query, searchTerm, productIDs)
	if !ok {
		logger.ERROR("Failed to get product IDs")
		return nil, false
	}


	return productIDs, true
}


func get(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

	searchTerm, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get search term")
		return false
	}

	rows, err := tx.QueryContext(ctx, query, searchTerm)
	if err != nil {
		logger.ERROR("Failed to execute the query. Reason: %s", err)
		return false
	}

	productIDs, ok := args[1].(*[]*string)
	if !ok {
		logger.ERROR("Failed to get product IDs")
		return false
	}

	for rows.Next() {
		var productID string
		err := rows.Scan(&productID)
		if err != nil {
			logger.ERROR("Failed to scan product ID")
			return false
		}
		*productIDs = append(*productIDs, &productID)
	}
	return true
}


func Add(logger *logger.Logger, tx *sql.Tx, searchTerm string, products *[]*product.Product) (ok bool) {

	searchTerm = strings.ToLower(searchTerm)
	
	lastFetched := int(time.Now().Unix())
	query := `INSERT INTO searches (search_term, product_id, fetch_count, last_fetch) VALUES ($1, $2, 1, $3)`

	ok = postgres.ExecuteContextChangeQuery(logger, tx, add, query, searchTerm, products, lastFetched)
	if !ok {
		logger.ERROR("Failed to add search term")
		return false
	}

	return true
}


func add(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {
	
	searchTerm, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get search term")
		return false
	}

	products, ok := args[1].(*[]*product.Product)
	if !ok {
		logger.ERROR("Failed to get product IDs")
		return false
	}

	lastFetched, ok := args[2].(int)
	if !ok {
		logger.ERROR("Failed to get last fetched")
		return false
	}

	for _, product := range *products {
		_, err := tx.ExecContext(ctx, query, searchTerm, product.ID, lastFetched)
		if err != nil {
			logger.ERROR("Failed to execute the query. Reason: %s", err)
			return false
		}
	}
	return true
}