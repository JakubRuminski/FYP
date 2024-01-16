package request

import (
	"net/http"

	"github.com/jakubruminski/FYP/go/api"
	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/token"
)

const pathToStaticBuild = "./static/build"


func HandleRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, requestID string) {
	token.CreateToken(logger, w, "User")

	fs := http.FileServer(http.Dir( pathToStaticBuild ))

	fs.ServeHTTP(w, r)
}

func HandleApiRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, requestID string) {
	logger.INFO("Request: %s", r.URL.Path)
	
	if !token.ValidToken(logger, r) {
		response.WriteResponse(logger, w, http.StatusUnauthorized, "application/json", "error", "Token not valid, please reload the page")
		return
	}

	_, ok := api.GetProducts(logger, r, w)
	if !ok {
		logger.ERROR("Error while getting products")
		response.WriteResponse(logger, w, http.StatusInternalServerError, "application/json", "error", "Problem getting products, try again or please return later.")
		return
	}
}