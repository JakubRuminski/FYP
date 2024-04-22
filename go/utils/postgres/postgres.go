package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


var (
	CONTEXT_TIMEOUT = 120
)


func credentialString(logger *logger.Logger) string {
	db_host := "DB_HOST"
	db_port := "DB_PORT"
	db_user := "POSTGRES_USER"
	db_password := "POSTGRES_PASSWORD"
	db_name := "POSTGRES_DB"

	ok := env.GetKeys(logger, &db_host, &db_port, &db_user, &db_password, &db_name)
	if !ok {
		logger.ERROR("Couldn't get environment variables")
	}

	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	return connectionString

}

func connectToDatabase(logger *logger.Logger) (db *sql.DB, ok bool) {
	psqlInfo := credentialString(logger)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.ERROR("DB Connection couldn't get established. Reason: %s", err)
		return nil, false
	}

	return db, true
}

func ExecuteCreateTableQuery(logger *logger.Logger, tableName, query string) (ok bool) {
	CONTEXT_TIMEOUT, ok := env.GetInt(logger, "CONTEXT_TIMEOUT")
	if !ok {
		return false
	}
	logger.DEBUG("Global context timout set to %d", CONTEXT_TIMEOUT)

	retry := true
	for retry {
		db, ok := connectToDatabase(logger)
		if !ok {
			logger.ERROR("Couldn't connect to the database")
			logger.DEBUG("Sleeping for 1 minute before retrying")
			time.Sleep(1 * time.Minute)
			continue
		}
		defer db.Close()

		_, err := db.Exec(query)
		if err != nil {
			logger.ERROR("Couldn't execute the query. Reason: %s", err)
			logger.DEBUG("Sleeping for 1 minute before retrying")
			time.Sleep(1 * time.Minute)
			continue
		}

		return true
	}

	logger.INFO("Successfully created table '%s'", tableName)
	return true
}

// REMEMBER TO COMMIT THE TRANSACTION
func createTransaction(logger *logger.Logger) (tx *sql.Tx, ok bool) {
	db, ok := connectToDatabase(logger)
	if !ok {
		logger.ERROR("Couldn't connect to the database")
		return nil, false
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		logger.ERROR("Couldn't start the transaction. Reason: %s", err)
		return nil, false
	}
	return tx, true
}

func ExecuteInTransaction(
	logger *logger.Logger,
	transactionFunction func(logger *logger.Logger, tx *sql.Tx, args ...interface{}) bool,
	transactionArgs ...interface{},
) bool {

	var ok bool
	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { return false }
	
	var tx *sql.Tx
	if db_available {
		tx, ok = createTransaction(logger)
		if !ok {
			logger.ERROR("Couldn't start the transaction.")
			return false
		}
		defer tx.Rollback()
	}

	ok = transactionFunction(logger, tx, transactionArgs...)
	if !ok {
		logger.ERROR("Transaction failed")
		return false
	}

	if db_available {
		err := tx.Commit()
		if err != nil {
			logger.ERROR("Couldn't commit the transaction. Reason: %s", err)
			return false
		}
	}

	return true
}

func ExecuteContextChangeQuery(
	logger *logger.Logger,
	tx *sql.Tx,
	changeFunction func(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (bool),
	query string,
	args ...interface{},
) (ok bool) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(CONTEXT_TIMEOUT)*time.Second)
	defer cancel()

	ok = changeFunction(logger, tx, ctx, query, args...)
	if !ok {
		logger.ERROR("Failed to execute the context change query")
		return false
	}

	return true
}

func ExecuteContextLookUpQuery(
	logger *logger.Logger,
	tx *sql.Tx,
	lookUpFunction func(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (bool),

	query string,
	args ...interface{},
) (ok bool) {

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Duration(CONTEXT_TIMEOUT)*time.Second)
	defer cancel()

	ok = lookUpFunction(logger, tx, ctx, query, args...)
	if !ok {
		logger.ERROR("Failed to execute the context look up query")
		return false
	}

	return true
}