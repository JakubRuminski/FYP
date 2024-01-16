package query

import (
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func Products(logger *logger.Logger, searchTerm, searchType string) (products *[]*product.Product, found, ok bool) {

	return nil, true, true

}