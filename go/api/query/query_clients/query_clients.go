package query_clients

import (
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
			client_id INT,
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
