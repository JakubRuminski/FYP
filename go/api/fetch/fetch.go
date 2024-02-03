package fetch

import (
	"sync"

	"github.com/jakubruminski/FYP/go/api/fetch/seller/tesco"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func Products(logger *logger.Logger, searchValue string) (products *[]*product.Product, ok bool) {

	var wg sync.WaitGroup

    products = &[]*product.Product{}

	wg.Add(1)
	go fetch(logger, tesco.Fetch, searchValue, &wg, products)

	wg.Wait()

	products, ok = product.Sort(logger, products)
	if !ok {
		logger.ERROR("Error while sorting products")
		return nil, false
	}

	return products, true

}


func fetch( logger *logger.Logger,
			fetchFunction func(logger *logger.Logger, searchValue string) (*[]*product.Product, bool),
			searchValue string,
			wg *sync.WaitGroup,
			products *[]*product.Product ) {

	defer wg.Done()

	fetchedProducts, ok := fetchFunction(logger, searchValue)
	if !ok {
		logger.ERROR("Error while fetching products")
		return
	}

	for _, product := range *fetchedProducts {
		*products = append(*products, product)
	}
}