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

	verbose, ok := env.GetBool(logger, "VERBOSE")
	if !ok { return false }

	logger.SetFlags(environment, verbose, clientID)

	file, ok := logger.InitRequestLogFile(clientID)
	if !ok {
		logger.ERROR("Error while creating log file")
		return false
	}
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

	logger.INFO("%s - %s", r.RemoteAddr, r.URL.Path)
	
	logger.DEBUG("%s", r.UserAgent())
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

	verbose, ok := env.GetBool(logger, "VERBOSE")
	if !ok { return false }

	logger.SetFlags(environment, verbose, clientID)

	logger.DEBUG("Client ID: %s", clientID)

	file, ok := logger.InitRequestLogFile(clientID)
	if !ok {
		logger.ERROR("Error while creating log file")
		return false
	}
	defer file.Close()

	_, ok = api.GetResponse(logger, r, w)
	if !ok {
		logger.ERROR("Error while getting products: %s", file.Name())
		return false
	}

	logger.DEBUG("Request handled. Client log file: %s", file.Name())

	return true
}
