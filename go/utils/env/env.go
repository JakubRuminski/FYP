package env


import (
	"fmt"
	"os"

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