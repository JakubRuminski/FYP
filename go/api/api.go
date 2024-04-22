package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"strings"

	"github.com/jakubruminski/FYP/go/api/fetch"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/api/query"

	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/http/response"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
	"github.com/jakubruminski/FYP/go/utils/token"
)

type Products struct {
	Results  *[]*product.Product               `json:"results"`
	Currency map[string]map[string]interface{} `json:"currency"`
}

func GetResponse(logger *logger.Logger, r *http.Request, w http.ResponseWriter) (jsonResponse []byte, ok bool) {
	// for {
	// 	// sleep for 1 second
	// 	time.Sleep(1 * time.Second)
	// }

	if r.URL.Path == "/api/search" {
		return getProductsHandler(logger, w, r)

	} else if r.URL.Path == "/api/add_item" {
		return addItemHandler(logger, w, r)

	} else if r.URL.Path == "/api/get_items" {
		return getItemsHandler(logger, w, r)

	} else if r.URL.Path == "/api/remove_item" {
		return removeItemHandler(logger, w, r)
	}

	logger.ERROR("Invalid request %s", r.URL.Path)
	return nil, false
}

func getProductsHandler(logger *logger.Logger, w http.ResponseWriter, r *http.Request) (jsonResponse []byte, ok bool) {
	searchTerm := parseSearchValue(r.FormValue("search_term"))
	searchTerm = strings.ToLower(searchTerm)

	products := &[]*product.Product{}

	ok = postgres.ExecuteInTransaction(logger, getProducts_DoInTransaction, products, searchTerm)
	if !ok && len(*products) == 0 {
		logger.ERROR("Failed to get products")
		return nil, false
	}
	if len(*products) == 0 {
		logger.ERROR("No products found")
		return nil, false
	}

	currency, ok := getCurrency(logger)
	if !ok {
		logger.ERROR("Failed to get currency")
		return nil, false
	}

	numberOfProductsForTesco := 0
	numberOfProductsForDunnes := 0
	numberOfProductsForSuperValu := 0
	for _, product := range *products {	
		if product.Seller == "Tesco" {
			numberOfProductsForTesco++
			
		} else if product.Seller == "Dunnes" {
			numberOfProductsForDunnes++

		} else if product.Seller == "SuperValu" {
			numberOfProductsForSuperValu++
		}
	}
	logger.INFO("Tesco: %d", numberOfProductsForTesco)
	logger.INFO("Dunnes: %d", numberOfProductsForDunnes)
	logger.INFO("SuperValu: %d", numberOfProductsForSuperValu)

	jsonResponse, err := json.Marshal(Products{Results: products, Currency: currency})
	if err != nil {
		logger.ERROR("Failed to marshal response: %s", err)
		return nil, false
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	logger.INFO("Client /logs/%s.txt searched for %s and got %d results", logger.ClientID, searchTerm, len(*products))

	return jsonResponse, true
}


func getProducts_DoInTransaction(logger *logger.Logger, tx *sql.Tx, args ...interface{}) bool {

    if len(args) != 2 {
        logger.ERROR("Expected 2 arguments, got %d", len(args))
        return false
    }

    products, ok := args[0].(*[]*product.Product)
    if !ok {
        logger.ERROR("Failed to get products")
        return false
    }

    searchTerm, ok := args[1].(string)
    if !ok {
		logger.ERROR("Failed to get search term")
        return false
    }

	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { return false }

	found := false
	expired := false
	if db_available {
		found, expired, ok := query.Products(logger, tx, products, searchTerm)

		if !ok {
			logger.ERROR("Failed to get products from database")
		}
		if !expired && found {
			logger.INFO("Products found in database")
			return true
		}
	}

	var oldProducts []*product.Product
	if expired {
		oldProducts = make([]*product.Product, len(*products))
		copy(oldProducts, *products)
	
		// Reset *products to an empty slice, without allocating new memory
		*products = (*products)[:0]
	
		logger.DEBUG_WARN("Products expired in database")
	}

	logger.DEBUG("products variable address: %p", products)
	logger.DEBUG("oldProducts variable address: %p", oldProducts)

	if !found {
		logger.DEBUG_WARN("No products matched in database, now web scraping...")
	}
	ok = fetch.Products(logger, products, searchTerm)
	if !ok {
		logger.ERROR("Failed to get products from web scraping")
		return false
	}

	if len(*products) == 0 {
		logger.ERROR("No products found")
		return true
	}
    if db_available {
		ok = query.AddProducts(logger, tx, searchTerm, &oldProducts, products)
		if !ok {
			logger.ERROR("Failed to add products to database")
			return false
		}

		ok = query.AddSearchTerm(logger, tx, searchTerm, products)
		if !ok {
			logger.ERROR("Failed to add search term to database")
			return false
		}
	}

	return true
}


// TODO: These should be fetched and not hardcoded.
func getCurrency(logger *logger.Logger) (Rates map[string]map[string]interface{}, ok bool) {
	Rates = map[string]map[string]interface{}{
		"Canada":     {"rate": 1.44, "symbol": "C$"},
		"India":      {"rate": 89.42, "symbol": "₹"},
		"Costa Rica": {"rate": 588.03, "symbol": "₡"},
		"Australia":  {"rate": 1.63, "symbol": "A$"},
		"UK":         {"rate": 0.86, "symbol": "£"},
		"Euro":       {"rate": 1.0, "symbol": "€"},
		"Poland":     {"rate": 4.46, "symbol": "zł"},
	}

	return Rates, true
}

// This function escapes html characters and replaces spaces with "%20"
func parseSearchValue(searchValue string) string {

	searchValue = strings.ToLower(searchValue)

	startQuote := ""
	endQuote := ""
	if strings.HasPrefix(searchValue, "\"") && strings.HasSuffix(searchValue, "\"") {
		searchValue = strings.Trim(searchValue, "\"")
		startQuote = "\""
		endQuote = "\""
	}

	searchValue = html.EscapeString(searchValue)

	// if there is a space, replace it with "%20"
	searchValue = strings.Replace(searchValue, " ", "%20", -1)

	searchValue = startQuote + searchValue + endQuote

	return searchValue
}


func addItemHandler(logger *logger.Logger, w http.ResponseWriter, r *http.Request) (jsonResponse []byte, ok bool) {
	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { return nil, false }

	if !db_available {
		message := "Sorry, could not add this product"
		response.WriteResponse(logger, w, http.StatusOK, "application/json", "message", message)
	}

	logger.INFO("Request: %s", r.URL.Path)

	clientID, ok := token.GetID(logger, r)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return nil, false
	}
     
	product, ok := product.ParseProduct(logger, r)
	if !ok {
		logger.ERROR("Failed to parse product")
		return nil, false
	}

	logger.DEBUG("Adding product to basket id: %d", product.ID)

	ok = postgres.ExecuteInTransaction(logger, AddItem_DoInTransaction, clientID, product)
	if !ok {
		logger.ERROR("Failed to add product to basket")
		return nil, false
	}

	message := fmt.Sprintf("Product added to basket: %s", product.Name)
	response.WriteResponse(logger, w, http.StatusOK, "application/json", "message", message)

	return nil, true
}


func AddItem_DoInTransaction(logger *logger.Logger, tx *sql.Tx, args ...interface{}) bool {
	if len(args) != 2 {
		logger.ERROR("Expected 2 arguments, got %d", len(args))
		return false
	}

	clientID, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return false
	}

	product, ok := args[1].(*product.Product)
	if !ok {
		logger.ERROR("Failed to get product")
		return false
	}

	return query.AddToBaskets(logger, tx, clientID, *product)
}


func getItemsHandler(logger *logger.Logger, w http.ResponseWriter, r *http.Request) (jsonResponse []byte, ok bool) {
	logger.INFO("Request: %s", r.URL.Path)

	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { return nil, false }

	if !db_available {
		message := "Sorry, cannot display products at this time."
		response.WriteResponse(logger, w, http.StatusOK, "application/json", "message", message)
		return nil, true
	}

	clientID, ok := token.GetID(logger, r)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return nil, false
	}

	products := &[]*product.Product{}
	ok = postgres.ExecuteInTransaction(logger, getItems_DoInTransaction, clientID, products)
	if !ok {
		logger.ERROR("Failed to get products")
		return nil, false
	}

	jsonResponse, err := json.Marshal(Products{Results: products})
	if err != nil {
		logger.ERROR("Failed to marshal response")
		return nil, false
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

	return jsonResponse, true
}


func getItems_DoInTransaction(logger *logger.Logger, tx *sql.Tx, args ...interface{}) bool {
	if len(args) != 2 {
		logger.ERROR("Expected 2 arguments, got %d", len(args))
		return false
	}

	clientID, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return false
	}

	products, ok := args[1].(*[]*product.Product)
	if !ok {
		logger.ERROR("Failed to get products")
		return false
	}

	return query.Baskets(logger, tx, clientID, products)
}

func removeItemHandler(logger *logger.Logger, w http.ResponseWriter, r *http.Request) (jsonResponse []byte, ok bool) {
	logger.INFO("Request: %s", r.URL.Path)

	db_available, ok := env.GetBool(logger, "DB_AVAILABLE")
	if !ok { return nil, false }

	if !db_available {
		message := "Sorry, could not remove this product"
		response.WriteResponse(logger, w, http.StatusOK, "application/json", "message", message)
	}

	clientID, ok := token.GetID(logger, r)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return nil, false
	}

	product, ok := product.ParseProduct(logger, r)
	if !ok {
		logger.ERROR("Failed to parse product")
		return nil, false
	}

	logger.DEBUG("Removing product from basket id: %d", product.ID)

	ok = postgres.ExecuteInTransaction(logger, RemoveItem_DoInTransaction, clientID, product)
	if !ok {
		logger.ERROR("Failed to remove product from basket")
		return nil, false
	}

	message := fmt.Sprintf("Product removed from basket: %s", product.Name)
	response.WriteResponse(logger, w, http.StatusOK, "application/json", "message", message)

	return nil, true
}

func RemoveItem_DoInTransaction(logger *logger.Logger, tx *sql.Tx, args ...interface{}) bool {
	if len(args) != 2 {
		logger.ERROR("Expected 2 arguments, got %d", len(args))
		return false
	}

	clientID, ok := args[0].(string)
	if !ok {
		logger.ERROR("Failed to get client ID")
		return false
	}

	product, ok := args[1].(*product.Product)
	if !ok {
		logger.ERROR("Failed to get product")
		return false
	}

	return query.RemoveFromBasket(logger, tx, clientID, *product)
}