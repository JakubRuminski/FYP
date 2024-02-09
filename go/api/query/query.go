package query

import (
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func INITIALISE_DATABASE(logger *logger.Logger) (ok bool) {

	return true
}

func AddProducts(logger *logger.Logger, products *[]*product.Product) (ok bool) {
	return false
}



func Products(logger *logger.Logger, searchTerm string) (products *[]*product.Product, found, ok bool) {

	return nil, false, true

}


func AddToBaskets(logger *logger.Logger, userID int, product product.Product) (ok bool) {

	return false

}


func Baskets(logger *logger.Logger, userID int) (products *[]*product.Product, ok bool) {

	return nil, false

}