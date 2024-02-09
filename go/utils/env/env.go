package env

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"

	"github.com/jakubruminski/FYP/go/utils/logger"
)

func LoadEnv() (ok bool) {
	err := godotenv.Load("/secrets/env/.env")
	if err != nil {
		fmt.Printf("Error loading .env file. Error: %v\n", err)
		return false
	}
	return true
}

func Get(logger *logger.Logger, key string) (value string, ok bool) {
	value, ok = os.LookupEnv(key)
	if !ok {
		logger.ERROR("Environment variable %s not set", key)
	}
	return
}

func GetInt(logger *logger.Logger, key string) (value int, ok bool) {
	v, ok := Get(logger, key)
	if !ok {
		return value, false
	}
	value, err := strconv.Atoi(v)
	if err != nil {
		logger.ERROR("Failed to convert %s to int. Reason: %s", key, err)
		return value, false
	}

	return value, true
}

func GetKeys(logger *logger.Logger, v ...*string) (ok bool) {
	exists := true
	for _, key := range v {
		*key, exists = Get(logger, *key)
		if !exists {
			ok = false
		}
	}
	ok = exists
	return ok
}