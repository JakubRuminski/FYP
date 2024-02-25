package query

import (
	"database/sql"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/api/query/query_clients"
	"github.com/jakubruminski/FYP/go/api/query/query_products"
	"github.com/jakubruminski/FYP/go/api/query/query_searchs"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func INITIALISE_DATABASE(logger *logger.Logger) (ok bool) {
    if !query_products.INIT(logger) {
        logger.ERROR("Failed to initialize products")
        return false
    }
    if !query_clients.INIT(logger) {
        logger.ERROR("Failed to initialize clients")
        return false
    }
    if !query_searchs.INIT(logger) {
        logger.ERROR("Failed to initialize searches")
        return false
    }

    return true
}

func Products(logger *logger.Logger, tx *sql.Tx, products *[]*product.Product, searchTerm string) (found, ok bool) {

    ProductIDs, ok := query_searchs.Get(logger, tx, searchTerm)
    if !ok {
        logger.ERROR("Failed to get product IDs")
        return false, false
    }

    if len(*ProductIDs) == 0 {
        logger.ERROR("No products found")
        return false, true
    }

    ok = query_products.Get(logger, tx, products, ProductIDs)
    if !ok {
        logger.ERROR("Failed to get products")
        return false, false
    }

	return true, true

}

func AddProducts(logger *logger.Logger, tx *sql.Tx, query string, productsToAdd *[]*product.Product) (ok bool) {

    if !query_products.Add(logger, tx, productsToAdd) {
        logger.ERROR("Failed to add products")
        return false
    }

	return true
}

func AddSearchTerm(logger *logger.Logger, tx *sql.Tx, searchTerm string, products *[]*product.Product) (ok bool) {
    
        if !query_searchs.Add(logger, tx, searchTerm, products) {
            logger.ERROR("Failed to add search term")
            return false
        }
    
        return true
    }

func AddToBaskets(logger *logger.Logger, clientID string, product product.Product) (ok bool) {

	return false

}

func Baskets(logger *logger.Logger, clientID string) (products *[]*product.Product, ok bool) {

	return nil, false

}
