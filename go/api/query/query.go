package query

import (
	"database/sql"
	"time"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/api/query/query_clients"
	"github.com/jakubruminski/FYP/go/api/query/query_products"
	"github.com/jakubruminski/FYP/go/api/query/query_searchs"
	"github.com/jakubruminski/FYP/go/utils/env"
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

func Products(logger *logger.Logger, tx *sql.Tx, products *[]*product.Product, searchTerm string) (found, expired, ok bool) {

    ProductIDs, ok := query_searchs.GetIDs(logger, tx, searchTerm)
    if !ok {
        logger.ERROR("Failed to get product IDs")
        return false, false, false
    }
    if len(*ProductIDs) == 0 {
        logger.DEBUG_WARN("No products found")
        return false, false, true
    }

    ok = query_products.Get(logger, tx, products, ProductIDs)
    if !ok {
        logger.ERROR("Failed to get products")
        return false, false, false
    }
    
    expiry_offset, ok := env.GetInt(logger, "SEARCH_EXPIRY_IN_DAYS")
    if !ok { return false, false, false }

    expiry_offset_seconds := expiry_offset * 24 * 60 * 60

    expiry, ok := query_searchs.GetExpiry(logger, tx, searchTerm)
    if !ok {
        logger.ERROR("Failed to get expiry")
        return false, false, false
    }

    nowTime := int(time.Now().Unix())
    if nowTime > (expiry + expiry_offset_seconds) {
        logger.DEBUG_WARN("Expiry time has passed")
        return false, true, true
    }

	return true, false, true

}

func AddProducts(logger *logger.Logger, tx *sql.Tx, query string, oldProducts, productsToAdd *[]*product.Product) (ok bool) {

    if !query_products.Add(logger, tx, oldProducts, productsToAdd) {
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

func AddToBaskets(logger *logger.Logger, tx *sql.Tx, clientID string, product product.Product) (ok bool) {

    if !query_clients.Add(logger, tx, clientID, product.ID) {
        logger.ERROR("Failed to add product to basket")
        return false
    }

    return true
}

func Baskets(logger *logger.Logger, tx *sql.Tx, clientID string, products *[]*product.Product) (ok bool) {

	if !query_clients.GetByID(logger, tx, clientID, products) {
        logger.ERROR("Failed to get products from basket")
        return false
    }

    return true
}

func RemoveFromBasket(logger *logger.Logger, tx *sql.Tx, clientID string, product product.Product) (ok bool) {
    
    if !query_clients.Remove(logger, tx, clientID, product.ID) {
        logger.ERROR("Failed to remove product from basket")
        return false
    }

    return true
}
