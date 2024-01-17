package query

import (
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func Products(logger *logger.Logger, searchTerm string) (products *[]*product.Product, found, ok bool) {

	return nil, false, true

}