package product

import (
	"encoding/json"
	"io"
	"net/http"
	"sort"

	"github.com/jakubruminski/FYP/go/utils/logger"
)

type Product struct {
	ID                   int64   `json:"id"`              // product id. Database specific
	Seller               string  `json:"seller"`
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

func ProductCreateQuery() (query string) {
	return `
	CREATE TABLE IF NOT EXISTS products (
		id                                  SERIAL PRIMARY KEY,
		seller                              VARCHAR(255),
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

func ProductInsertQuery() (query string) {
	query = `
    INSERT INTO products
    (seller, name, currency, price, price_per_unit, discount_price, discount_price_per_unit, discount_price_in_words, unit_type, url, img_url)
    VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id
    `
	return query
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
	var result struct {
        Result *Product `json:"result"`
    }

    body, err := io.ReadAll(r.Body)
    if err != nil {
        logger.ERROR("Failed to read request body. Reason: %s", err)
        return nil, false
    }

    if err := json.Unmarshal(body, &result); err != nil {
        logger.ERROR("Failed to decode product. Reason: %s", err)
        return nil, false
    }

    if result.Result == nil {
        logger.ERROR("Missing 'result' field in JSON.")
        return nil, false
    }

    product = result.Result

    logger.DEBUG("p.ID: %d", product.ID)
    logger.DEBUG("p.Name: %s", product.Name)
    logger.DEBUG("p.Seller: %s", product.Seller)
    logger.DEBUG("p.Price: %f", product.Price)
    logger.DEBUG("p.PricePerUnit: %f", product.PricePerUnit)
    logger.DEBUG("p.DiscountPrice: %f", product.DiscountPrice)
    logger.DEBUG("p.DiscountPricePerUnit: %f", product.DiscountPricePerUnit)
    logger.DEBUG("p.DiscountPriceInWords: %s", product.DiscountPriceInWords)
    logger.DEBUG("p.UnitType: %s", product.UnitType)
    logger.DEBUG("p.URL: %s", product.URL)
    logger.DEBUG("p.ImgURL: %s", product.ImgURL)

    return product, true
}


func initProduct(logger *logger.Logger, currency, seller, name string, price, pricePerUnit, discountPrice, discountPricePerUnit float64, DiscountPriceInWords, UnitType, url, imgURL string) (product *Product) {
	product = new(Product)
	product.ID = -1              // This is a placeholder for the product id. It is replaced by actual next available id in the database.
	product.Seller = seller
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


func Sort(logger *logger.Logger, products *[]*Product) (ok bool) {
	if len(*products) == 0 {
		logger.DEBUG_WARN("No products to sort")
		return true
	}

	sort.SliceStable((*products), func(i, j int) bool {
		Product_i_PricePerUnit := (*products)[i].PricePerUnit
		Product_i_DiscountPricePerUnit := (*products)[i].DiscountPricePerUnit
		Product_j_PricePerUnit := (*products)[j].PricePerUnit
		Product_j_DiscountPricePerUnit := (*products)[j].DiscountPricePerUnit

		bothHaveDiscount := (Product_i_DiscountPricePerUnit != 0.0 && Product_j_DiscountPricePerUnit != 0.0)
		firstProductHasDiscount := Product_i_DiscountPricePerUnit != 0.0
		secondProductHasDiscount := Product_j_DiscountPricePerUnit != 0.0

		if bothHaveDiscount {
			return Product_i_DiscountPricePerUnit < Product_j_DiscountPricePerUnit
		}
		if firstProductHasDiscount {
			return Product_i_DiscountPricePerUnit < Product_j_PricePerUnit
		}
		if secondProductHasDiscount {
			return Product_i_PricePerUnit < Product_j_DiscountPricePerUnit
		}

		if Product_i_PricePerUnit == Product_j_PricePerUnit {
			return false
		}
		
		return Product_i_PricePerUnit < Product_j_PricePerUnit
	})

	for i, product := range *products {
		if product.DiscountPricePerUnit != 0.0 {
			logger.DEBUG("Product %d: %f", i, product.DiscountPricePerUnit)
		} else {
			logger.DEBUG("Product PricePerUnit: %f", product.PricePerUnit)
		}
	}

	return true
}
