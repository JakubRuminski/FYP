package query_searchs

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "searches"

type SearchTerm struct {
	SearchTerm          string     `json:"search_term"`
	ProductID		    int        `json:"product_id"`

	FetchCount          int        `json:"fetch_count"`
	LastFetch           int        `json:"last_fetch"`
	Expiry 			    int        `json:"expiry"`
}

func INIT(logger *logger.Logger) (ok bool) {
	
	query := `
		CREATE TABLE IF NOT EXISTS searches
		(
			search_term      VARCHAR(255),
			product_id       INT,
			fetch_count      INT,
			last_fetch       INT,
			expiry           INT
		)
	`

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the searches table")
		return false
	}

	return true	
}

func GetIDs(logger *logger.Logger, tx *sql.Tx, searchTerm string) (productIDs *[]*int64, ok bool) {
	searchTerm = strings.ToLower(searchTerm)

	query := `SELECT product_id FROM searches WHERE search_term = $1`

	productIDs = &[]*int64{}
	ok = postgres.ExecuteContextLookUpQuery(logger, tx, getIDs, query, searchTerm, productIDs)
	if !ok {
		logger.ERROR("Failed to get product IDs")
		return nil, false
	}


	return productIDs, true
}


func getIDs(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

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

	productIDs, ok := args[1].(*[]*int64)
	if !ok {
		logger.ERROR("Failed to get product IDs")
		return false
	}

	for rows.Next() {
		productID := new(int64)
		err := rows.Scan(productID)
		if err != nil {
			logger.ERROR("Failed to scan product ID")
			return false
		}
		*productIDs = append(*productIDs, productID)
	}
	return true
}


func GetExpiry(logger *logger.Logger, tx *sql.Tx, searchTerm string) (expiry int, ok bool) {
	searchTerm = strings.ToLower(searchTerm)

	query := `SELECT last_fetch FROM searches WHERE search_term = $1`

	ok = postgres.ExecuteContextLookUpQuery(logger, tx, getExpiry, query, searchTerm, &expiry)
	if !ok {
		logger.ERROR("Failed to get expiry")
		return 0, false
	}

	return expiry, true
}


func getExpiry(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {
	
	searchTerm, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get search term")
		return false
	}

	expiry, ok := args[1].(*int)
	if !ok {
		logger.ERROR("Failed to get expiry")
		return false
	}

	err := tx.QueryRowContext(ctx, query, searchTerm).Scan(expiry)
	if err != nil {
		logger.ERROR("Failed to execute the query. Reason: %s", err)
		return false
	}
	return true
}


func Add(logger *logger.Logger, tx *sql.Tx, searchTerm string, products *[]*product.Product) (ok bool) {

	expiry, ok := env.GetInt(logger, "SEARCH_EXPIRY_IN_DAYS")
	if !ok {
		logger.ERROR("Failed to get search expiry")
		return false
	}
	searchTerm = strings.ToLower(searchTerm)
	
	now := int(time.Now().Unix()) 
	lastFetched := now
	expiry = now + (expiry * 24 * 60 * 60)
	query := `INSERT INTO searches (search_term, product_id, fetch_count, last_fetch, expiry) VALUES ($1, $2, 1, $3, $4)`

	ok = postgres.ExecuteContextChangeQuery(logger, tx, add, query, searchTerm, products, lastFetched, expiry)
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

	expiry, ok := args[3].(int)
	if !ok {
		logger.ERROR("Failed to get expiry")
		return false
	}

	for _, product := range *products {
		_, err := tx.ExecContext(ctx, query, searchTerm, product.ID, lastFetched, expiry)
		if err != nil {
			logger.ERROR("Failed to execute the query. Reason: %s", err)
			return false
		}
	}
	return true
}