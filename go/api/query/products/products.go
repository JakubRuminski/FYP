package products

import (
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


type SearchTerms struct {
	SearchTerm          string     `json:"search_term"`

	ProductID		    int        `json:"product_id"`
}




type Product struct {
	FetchCount          int        `json:"fetch_count"`

	product.Product
}

func getProductsColumnsAndTypes(logger *logger.Logger) (columnsAndTypes map[string]string) {

    columnsAndTypes = map[string]string{
        "fetch_count":               "INT",
        "seller":                    "VARCHAR(255)",
        "id":                        "PRIMARY KEY AUTOINCREMENT",
        "name":                      "VARCHAR(255)",
		"currency":                  "VARCHAR(10)",
		"price":                     "DOUBLE PRECISION",
        "price_per_unit":            "DOUBLE PRECISION",
        "discount_price":            "DOUBLE PRECISION",
        "discount_price_per_unit":   "DOUBLE PRECISION",
        "discount_price_in_words":   "VARCHAR(255)",
        "unit_type":                 "VARCHAR(30)",
		"url":                       "VARCHAR(255)",
        "img_url":                   "VARCHAR(255)",
    }

    return columnsAndTypes
}


func INIT_PRODUCTS_TABLES(logger *logger.Logger) (ok bool) {

	columnsAndTypes := getProductsColumnsAndTypes(logger)

	query := "CREATE TABLE IF NOT EXISTS products ("
	for column, columnType := range columnsAndTypes {
		query += column + " " + columnType + ", "
	}

	// TODO: CONTINUE HERE


	return true
	
}