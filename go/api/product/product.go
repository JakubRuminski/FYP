package product

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/jakubruminski/FYP/go/utils/logger"
)

type Product struct {
	Seller               string  `json:"seller"`
	ID                   string  `json:"id"`              // product id. Database specific
	Name                 string  `json:"name"`

	Currency			 string  `json:"currency"`
	Price                float64 `json:"price"`
	PricePerUnit         float64 `json:"price_per_unit"`
	DiscountPrice        float64 `json:"discount_price"`
	DiscountPricePerUnit float64 `json:"discount_price_per_unit"`
	DiscountPriceInWords string  `json:"discount_price_in_words"`
	UnitType             string  `json:"unit_type"`

	URL                  string  `json:"url"`
	ImgURL               string  `json:"img_url"`
}

func ProductQuery() (query string) {
	return `
	CREATE TABLE IF NOT EXISTS products (
		seller                              VARCHAR(255),
		id                                  SERIAL PRIMARY KEY,
		name                                VARCHAR(255),
		currency                            VARCHAR(10),
		price                               DOUBLE PRECISION,
		price_per_unit                      DOUBLE PRECISION,
		discount_price                      DOUBLE PRECISION,
		discount_price_per_unit             DOUBLE PRECISION,
		discount_price_in_words             VARCHAR(255),
		unit_type                           VARCHAR(30),
		url                                 VARCHAR(255),
		img_url                             VARCHAR(255)
	)	
	`
}


func NewProduct(logger                     *logger.Logger, 
	            seller                     string,
				name                       string,
				currency                   string,
				price                      float64,
				pricePerUnit               float64,
				discountPrice              float64,
				discountPricePerUnit       float64,
				pricePerUnitUnitType       string, 
				DiscountPriceInWords       string,
				url                        string,
				imgURL                     string) (product *Product, ok bool) {


	product = initProduct(logger, currency, seller, name, price, pricePerUnit, discountPrice, discountPricePerUnit, DiscountPriceInWords, pricePerUnitUnitType, url, imgURL)

	return product, true
}


func ParseProduct(logger *logger.Logger, r *http.Request) (product *Product, ok bool) {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	if err := decoder.Decode(&product); err != nil {
		logger.ERROR("Failed to decode product. Reason: %s", err)
		return nil, false
	}

	return product, true
}


func initProduct(logger *logger.Logger, currency, seller, name string, price, pricePerUnit, discountPrice, discountPricePerUnit float64, DiscountPriceInWords, UnitType, url, imgURL string) (product *Product) {
	product = new(Product)
	product.Seller = seller
	product.ID = "ID"      // This is a placeholder for the product id. It is replaced by actual next available id in the database.
	product.Name = name
	
	product.Currency = currency
	product.Price = price
	product.PricePerUnit = pricePerUnit
	product.DiscountPrice = discountPrice
	product.DiscountPricePerUnit = discountPricePerUnit
	product.DiscountPriceInWords = DiscountPriceInWords
	product.UnitType = UnitType

	product.URL = url
	product.ImgURL = imgURL
	return product
}


func Sort(logger *logger.Logger, products *[]*Product) (sortedProducts *[]*Product, ok bool) {
	if len(*products) == 0 {
		logger.DEBUG_WARN("No products to sort")
		return nil, false
	}

	// Use sort.SliceStable to sort the slice if you want to preserve the original order among equal elements.
	// Otherwise, you can use sort.Slice for potentially faster sorting without this guarantee.
	sort.SliceStable((*products), func(i, j int) bool {
		Product_i_PricePerUnit := (*products)[i].PricePerUnit
		Product_i_DiscountPrice := (*products)[i].DiscountPrice
		Product_j_PricePerUnit := (*products)[j].PricePerUnit
		Product_j_DiscountPrice := (*products)[j].DiscountPrice

		bothHaveDiscount := Product_i_DiscountPrice != 0 && Product_j_DiscountPrice != 0
		firstProductHasDiscount := Product_i_DiscountPrice != 0
		secondProductHasDiscount := Product_j_DiscountPrice != 0

		if bothHaveDiscount {
			return Product_i_DiscountPrice < Product_j_DiscountPrice
		}
		if firstProductHasDiscount {
			return Product_i_DiscountPrice < Product_j_PricePerUnit
		}
		if secondProductHasDiscount {
			return Product_i_PricePerUnit < Product_j_DiscountPrice
		}

		if Product_i_PricePerUnit == Product_j_PricePerUnit {
			return false
		}

		return Product_i_PricePerUnit < Product_j_PricePerUnit
	})

	return products, true
}
