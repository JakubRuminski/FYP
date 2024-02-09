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
	db_user      := "DB_USER"
	db_password  := "DB_PASSWORD"
	db_name      := "DB_NAME"

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

