package api

import (
	"encoding/json"
	"html"
	"net/http"

	"github.com/jakubruminski/FYP/go/api/fetch"
	"github.com/jakubruminski/FYP/go/api/query"

	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

type Products struct {
	Results  []byte `json:"results"`
	Currency []byte `json:"currency"`
}


func GetProducts(logger *logger.Logger, r *http.Request, w http.ResponseWriter) (jsonResponse []byte, ok bool) {

	searchTerm := html.EscapeString( r.FormValue("search_term") )
	searchType := html.EscapeString( r.FormValue("search_type") )

    products, ok := getProducts(logger, searchTerm, searchType)
	if !ok {
		response.WriteResponse( logger, w, http.StatusInternalServerError, "application/json", "error", "Failed to get products" )
		return nil, false
	}

	currency, ok := getCurrency(logger)
	if !ok {
		response.WriteResponse( logger, w, http.StatusInternalServerError, "application/json", "error", "Failed to get currency" )
		return nil, false
	}

	jsonResponse, err := json.Marshal(Products{ Results:  products, Currency: currency })
	if err != nil {
		logger.ERROR("Failed to marshal products. Reason: %s", err)
		response.WriteResponse( logger, w, http.StatusInternalServerError, "application/json", "error", "Failed to marshal products" )
		return nil, false
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return jsonResponse, true
}


func getProducts(logger *logger.Logger, searchTerm, searchType string) (jsonResponse []byte, ok bool) {
	products, found, ok := query.Products(logger, searchTerm, searchType)
	if !ok {
		logger.ERROR("Failed to get products from database")
		return nil, false
	}
	
	if !found {
		logger.ERROR("No products matched in database, now web scraping...")
		products, ok = fetch.Products(logger, searchTerm, searchType)

		if !ok {
			logger.ERROR("Failed to get products from web scraping")
			return nil, false
		}
	}

	jsonResponse, err := json.Marshal(products)
	if err != nil {
		logger.ERROR("Failed to marshal products. Reason: %s", err)
		return nil, false
	}

	return jsonResponse, true
}


func getCurrency(logger *logger.Logger) (jsonResponse []byte, ok bool) {
	var Rates = map[string]map[string]interface{} {
		"Canada":     {"rate": 1.44, "symbol": "C$"},
		"India":      {"rate": 89.42, "symbol": "₹"},
		"Costa Rica": {"rate": 588.03, "symbol": "₡"},
		"Australia":  {"rate": 1.63, "symbol": "A$"},
		"UK":         {"rate": 0.86, "symbol": "£"},
		"Euro":       {"rate": 1.0, "symbol": "€"},
		"Poland":     {"rate": 4.46, "symbol": "zł"},
	}

	jsonResponse, err := json.Marshal(Rates)
	if err != nil {
		logger.ERROR("Failed to marshal currency. Reason: %s", err)
		return nil, false
	}

	return jsonResponse, true
}