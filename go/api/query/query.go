package query

import (
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/api/query/clients"
	"github.com/jakubruminski/FYP/go/api/query/products"
	"github.com/jakubruminski/FYP/go/api/query/searchs"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func INITIALISE_DATABASE(logger *logger.Logger) (ok bool) {
    if !products.INIT(logger) {
        logger.ERROR("Failed to initialize products")
        return false
    }
    if !clients.INIT(logger) {
        logger.ERROR("Failed to initialize clients")
        return false
    }
    if !searchs.INIT(logger) {
        logger.ERROR("Failed to initialize searches")
        return false
    }

    return true
}

func AddProducts(logger *logger.Logger, query string, productsToAdd *[]*product.Product) (ok bool) {

    if !products.Add(logger, productsToAdd) {
        logger.ERROR("Failed to add products")
        return false
    
    }

	return false
}

func Products(logger *logger.Logger, searchTerm string) (products *[]*product.Product, found, ok bool) {

	return nil, false, true

}

func AddToBaskets(logger *logger.Logger, clientID int, product product.Product) (ok bool) {

	return false

}

func Baskets(logger *logger.Logger, clientID int) (products *[]*product.Product, ok bool) {

	return nil, false

}
