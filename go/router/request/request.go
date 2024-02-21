package request

import (
	"net/http"

	"github.com/jakubruminski/FYP/go/api"
	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/token"
)

const pathToStaticBuild = "./static/build"

func HandleRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, clientID string) {
	_, ok := token.ValidToken(logger, r)
	if !ok {
		ok = token.CreateToken(logger, w, clientID)
	}

	if !ok {
		logger.ERROR("Error while creating token")
		response.WriteResponse(logger, w, http.StatusInternalServerError, "application/json", "error", "Something went wrong, please try again.")
		return
	}

	fs := http.FileServer(http.Dir(pathToStaticBuild))

	fs.ServeHTTP(w, r)
	logger.DEBUG("Static files served.")
}


func HandleApiRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, clientID string) {
	ok := handleApiRequest(w, r, logger, clientID)

	if !ok {
		response.WriteResponse(logger, w, http.StatusInternalServerError, "application/json", "error", "Something went wrong, please try again.")
	}
}

func handleApiRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, clientID string) (ok bool) {
	logger.INFO("Request: %s", r.URL.Path)
	_, ok = token.ValidToken(logger, r)
	if !ok {
		logger.ERROR("Error while validating token")
		ok = token.CreateToken(logger, w, clientID)
	}
	if !ok {
		logger.ERROR("Error while creating token")
		return false
	}

	_, ok = api.GetResponse(logger, r, w)
	if !ok {
		logger.ERROR("Error while getting products")
		return false
	}

	return true
}
