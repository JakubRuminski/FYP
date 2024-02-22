package request

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/jakubruminski/FYP/go/api"
	
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/token"
)

const pathToStaticBuild = "./static/build"

func HandleRequest(w http.ResponseWriter, r *http.Request) {

	logger := &logger.Logger{}

	ok := handleRequest(w, r, logger)
	if !ok {
		response.WriteResponse(logger, w, http.StatusInternalServerError, "application/json", "error", "Something went wrong, please try again.")
	}

}


func handleRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger) (ok bool) {
    
	clientID, ok := token.GetID(logger, r)
	if !ok {
		clientID = uuid.New().String()
		ok = token.CreateToken(logger, w, clientID)
		if !ok {
			logger.ERROR("Error while creating token")
			return false
		}
	}

	environment, ok := env.Get(logger, "ENVIRONMENT")
	if !ok { return false }
	logger.SetEnvironment(environment)

	file := logger.InitRequestLogFile(clientID)
	defer file.Close()

	fs := http.FileServer(http.Dir(pathToStaticBuild))

	fs.ServeHTTP(w, r)
	logger.DEBUG("Static files served.")

	return true
}


func HandleApiRequest(w http.ResponseWriter, r *http.Request) {
	logger := &logger.Logger{}
	ok := handleApiRequest(w, r, logger)

	if !ok {
		response.WriteResponse(logger, w, http.StatusInternalServerError, "application/json", "error", "Something went wrong, please try again.")
	}
}

func handleApiRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger) (ok bool) {
	logger.INFO("Request: %s", r.URL.Path)

	clientID, ok := token.GetID(logger, r)
	if !ok {
		clientID = uuid.New().String()
		ok = token.CreateToken(logger, w, clientID)
		if !ok {
			logger.ERROR("Error while creating token")
			return false
		}
	}

	environment, ok := env.Get(logger, "ENVIRONMENT")
	if !ok { return false }
	logger.SetEnvironment(environment)

	logger.INFO("Client ID: %s", clientID)

	file := logger.InitRequestLogFile(clientID)
	defer file.Close()

	_, ok = api.GetResponse(logger, r, w)
	if !ok {
		logger.ERROR("Error while getting products")
		return false
	}

	return true
}
