package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func credentialString(logger *logger.Logger) string {
	db_host      := "DB_HOST"
	db_port      := "DB_PORT"
	db_user      := "POSTGRES_USER"
	db_password  := "POSTGRES_PASSWORD"
	db_name      := "POSTGRES_DB"

	ok := env.GetKeys(logger, &db_host, &db_port, &db_user, &db_password, &db_name)
	if !ok {
		logger.ERROR("Couldn't get environment variables")
	}
	
	connectionString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", db_host, db_port, db_user, db_password, db_name)
	return connectionString
	
}

func connectToDatabase( logger *logger.Logger ) (db *sql.DB, ok bool) {
	psqlInfo := credentialString(logger)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		logger.ERROR("DB Connection couldn't get established. Reason: %s", err)
		return nil, false
	}

	return db, true
}


func ExecuteCreateTableQuery(logger *logger.Logger, tableName, query string) (ok bool) {
	db, ok := connectToDatabase(logger)
	if !ok {
		logger.ERROR("Couldn't connect to the database")
		return false
	}
	defer db.Close()

	_, err := db.Exec(query)
	if err != nil {
		logger.ERROR("Couldn't execute the query. Reason: %s", err)
		return false
	}

	logger.DEBUG("Successfully created table '%s'", tableName)
	return true
}


func ExecuteInTransaction(logger *logger.Logger, query string, values ...interface{}) (ok bool) {
	db, ok := connectToDatabase(logger)
	if !ok {
		logger.ERROR("Couldn't connect to the database")
		return false
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		logger.ERROR("Couldn't start the transaction. Reason: %s", err)
		return false
	}

	_, err = tx.Exec(query, values...)
	if err != nil {
		logger.ERROR("Couldn't execute the query. Reason: %s", err)
		err = tx.Rollback()
		if err != nil {
			logger.ERROR("Couldn't rollback the transaction. Reason: %s", err)
		}
		return false
	}

	err = tx.Commit()
	if err != nil {
		logger.ERROR("Couldn't commit the transaction. Reason: %s", err)
		return false
	}

	return true
}
