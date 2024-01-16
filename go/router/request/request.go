package request

import (
	"net/http"
	"os"

	"github.com/jakubruminski/FYP/go/api"
	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/token"
)

const pathToStaticBuild = "./static/build"


func HandleRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, requestID string) {
	token.CreateToken(logger, w, "User")
	
	// if './static/build' doesn't exist return 404
	if _, err := os.Stat(pathToStaticBuild); err != nil {
		response.WriteResponse(logger,
			                   w,
							   http.StatusInternalServerError,
							   "text/html",
							   "error",
							   "We incountered a problem, please try again later.")

    }
	
	fs := http.FileServer(http.Dir("./static/build"))

	fs.ServeHTTP(w, r)
}

func HandleApiRequest(w http.ResponseWriter, r *http.Request, logger *logger.Logger, requestID string) {
	if !token.ValidToken(logger, r) {
		logger.ERROR("Token not valid, ignoring API request")
		w.WriteHeader(http.StatusUnauthorized)
		w.Header().Set("Content-Type", "application/json")
		error := `{
			"error": "Token not valid, please reload the page"
		}`
		w.Write([]byte(error))
		return
	}

	jsonResponse, ok := api.GetProducts(logger, r, w)
	if !ok {
		logger.ERROR("Error while getting products")
		w.WriteHeader(http.StatusInternalServerError)
		w.Header().Set("Content-Type", "application/json")
		error := `{
			"error": "Problem getting products, try again or please return later."
		}`
		w.Write([]byte(error))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}