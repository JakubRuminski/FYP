package product

import (

	p "github.com/jakubruminski/FYP/go/utils/parse/price"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

type Product struct {
	Seller               string  `json:"seller"`
	ID                   string  `json:"id"`              // product id. Database specific
	Name                 string  `json:"name"`
	Currency			 string  `json:"currency"`
	Price                float64 `json:"price"`
	PricePerUnit         float64 `json:"sub_price"`
	DiscountPrice        float64 `json:"discount_price"`
	DiscountPricePerUnit float64 `json:"discount_price_per_unit"`
	URL                  string  `json:"url"`
	ImgURL               string  `json:"img_url"`
}

func NewProduct(logger *logger.Logger, seller, id, name, price, pricePerUnit, discountPrice, url, imgURL string) (product *Product, ok bool) {

	currency, priceFloat, ok := p.Float(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price %s", price)
		return nil, false
	}

	currency, pricePerUnitFloat, pricePerUnit, ok := p.FloatPerUnit(logger, pricePerUnit)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price %s", pricePerUnit)
		return nil, false
	}

	currency, discountPriceFloat, ok := p.Float(logger, discountPrice)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price %s", discountPrice)
		return nil, false
	}

	discountPricePerUnit := ( (priceFloat / discountPriceFloat) * pricePerUnitFloat )

	product = initProduct(logger, currency, seller, id, name, priceFloat, pricePerUnitFloat, discountPriceFloat, discountPricePerUnit, url, imgURL)

	return product, true
}

func initProduct(logger *logger.Logger, currency, seller, id, name string, price, pricePerUnit, discountPrice, discountPricePerUnit float64, url, imgURL string) (product *Product) {
	product = new(Product)
	product.Seller = seller
	product.ID = id
	product.Name = name
	
	product.Currency = currency
	product.Price = price
	product.PricePerUnit = pricePerUnit
	product.DiscountPrice = discountPrice
	product.DiscountPricePerUnit = discountPricePerUnit

	product.URL = url
	product.ImgURL = imgURL
	return product
}