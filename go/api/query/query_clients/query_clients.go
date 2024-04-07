package query_clients

import (
	"context"
	"database/sql"
	"time"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/api/query/query_products"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "clients"

type Client struct {
	ClientID      int  `json:"client_id"`      // Client's unique ID
	LastFetch     int  `json:"last_fetch"`     // Timestamp
	ProductID     int  `json:"product_id"`     // Reference to the product in the products table
	ProductExists bool `json:"product_exists"` // This is a flag allowing to check if the product exists
}

func INIT(logger *logger.Logger) (ok bool) {

	query := `
		CREATE TABLE IF NOT EXISTS clients
		(
			id SERIAL PRIMARY KEY,
			client_id  VARCHAR(255),
			last_fetch INT,
			product_id INT,
			product_exists BOOLEAN
		)
	`

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the clients table")
		return false
	}

	return true
}


func Add(logger *logger.Logger, tx *sql.Tx, clientID string, productID int64) (ok bool) {

	lastFetched := time.Now().Unix()
	productExists := true

	query := `
		INSERT INTO clients (client_id, last_fetch, product_id, product_exists)
		VALUES ($1, $2, $3, $4)
	`

	ok = postgres.ExecuteContextChangeQuery(logger, tx, add, query, clientID, lastFetched, productID, productExists)
	if !ok {
		logger.ERROR("Failed to add client")
		return false
	}

	return true
}

func add(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

	clientID, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return false
	}

	lastFetched, ok := args[1].(int64)
	if !ok {
		logger.ERROR("Failed to get last fetched")
		return false
	}

	productID, ok := args[2].(int64)
	if !ok {
		logger.ERROR("Failed to get product ID")
		return false
	}

	productExists, ok := args[3].(bool)
	if !ok {
		logger.ERROR("Failed to get product exists")
		return false
	}

	_, err := tx.ExecContext(ctx, query, clientID, lastFetched, productID, productExists)
	if err != nil {
		logger.ERROR("Failed to add client: %s", err)
		return false
	}

	return true
}

func GetByID(logger *logger.Logger, tx *sql.Tx, clientID string, products *[]*product.Product) (ok bool) {
	
	query := `
		SELECT product_id FROM clients
		WHERE client_id = $1
	`

	productIDs := &[]*int64{}
	ok = postgres.ExecuteContextLookUpQuery(logger, tx, getByID, query, clientID, productIDs)
	if !ok {
		logger.ERROR("Failed to get clients")
		return false
	}


	ok = query_products.Get(logger, tx, products, productIDs)
	if !ok {
		logger.ERROR("Failed to get products")
		return false
	}

	return true
}

func getByID(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {
	
	clientID, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return false
	}

	productIDs, ok := args[1].(*[]*int64)
	if !ok {
		logger.ERROR("Failed to get products")
		return false
	}

	rows, err := tx.QueryContext(ctx, query, clientID)
	if err != nil {
		logger.ERROR("Failed to get clients: %s", err)
		return false
	}
	defer rows.Close()

	for rows.Next() {
		productID := new(int64)
		err := rows.Scan(&productID)
		if err != nil {
			logger.ERROR("Failed to scan client: %s", err)
			return false
		}
		*productIDs = append(*productIDs, productID)
	}

	return true
}