package response

import (
	"fmt"
	"net/http"

	"github.com/jakubruminski/FYP/go/utils/logger"
)


// WriteResponse writes a response to the client
//
// Parameters:
//  logger:                 The logger
//  w:                      The response writer
//  statusCode:             The status code to write.      e.g. http.StatusBadRequest
//  contentType:            The content type to write      e.g. "application/json"
//  responseToClient_Type:  The type of response to write  e.g. "data", "message" or "error"
//  responseToClient:       The response to write          e.g. "Invalid search term"
//
func WriteResponse(logger *logger.Logger, w http.ResponseWriter, statusCode int, contentType string, responseToClient_Type, responseToClient string) {
	logger.DEBUG("Writing response to client: %s", responseToClient)
	
	if w.Header().Get("Status") == "" {
		w.WriteHeader(statusCode)
	}
	if w.Header().Get("Content-Type") == "" {
	    w.Header().Set("Content-Type", contentType)
	}

	if responseToClient_Type != "data" && responseToClient_Type != "message" && responseToClient_Type != "error" {
		logger.ERROR("Invalid response type: %s", responseToClient_Type)
		responseToClient_Type = "error"
		return
	}

	responseToClient = fmt.Sprintf(`{"%s": "%s"}`, responseToClient_Type, responseToClient)

	w.Write([]byte(responseToClient))
}