package products

import (
	"strings"

	"github.com/jakubruminski/FYP/go/api/product"

	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "products"

func INIT(logger *logger.Logger) (ok bool) {
    
	query := product.ProductQuery()

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the products table")
		return false
	}

	return true	
}


func Add(logger *logger.Logger, products *[]*product.Product) bool {
    // Check if there are products to insert
    if len(*products) == 0 {
        logger.INFO("No products to add")
        return true
    }

    // Start building the batch insert query
    query := `
    INSERT INTO products 
    (seller, id, name, currency, price, price_per_unit, discount_price, discount_price_per_unit, discount_price_in_words, unit_type, url, img_url)
    VALUES `

    // Placeholder slice and values slice
    var placeholders []string
    var values []interface{}

    // Loop through products to prepare the placeholders and values
    for _, product := range *products {
        placeholders = append(placeholders, "(1$, 2$, 3$, 4$, 5$, 6$, 7$, 8$, 9$, 10$, 11$, 12$)")
        values = append(values,
            product.Seller,
            product.ID,
            product.Name,
            product.Currency,
            product.Price,
            product.PricePerUnit,
            product.DiscountPrice,
            product.DiscountPricePerUnit,
            product.DiscountPriceInWords,
            product.UnitType,
            product.URL,
            product.ImgURL,
        )
    }

    // Join all placeholders for the query and add to the main query string
    query += strings.Join(placeholders, ",")

    ok := postgres.ExecuteInTransaction(logger, query, values...)
    if !ok {
        logger.ERROR("Failed to batch insert products")
        return false
    }

    return true
}