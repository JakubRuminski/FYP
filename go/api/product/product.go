package product

import (
	"encoding/json"
	"net/http"
	"sort"

	"github.com/jakubruminski/FYP/go/utils/logger"
	p "github.com/jakubruminski/FYP/go/utils/parse/price"
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
	DiscountPriceInWords string  `json:"discount_price_in_words"`
	UnitType             string  `json:"unit_type"`

	URL                  string  `json:"url"`
	ImgURL               string  `json:"img_url"`
}

func NewProduct(logger *logger.Logger, seller, id, name, price, pricePerUnit, discountPrice, DiscountPriceInWords, url, imgURL string) (product *Product, ok bool) {

	currency, priceFloat, ok := p.Float(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price %s", price)
		return nil, false
	}

	_, pricePerUnitFloat, unitType, ok := p.FloatPerUnit(logger, pricePerUnit)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price %s", pricePerUnit)
		return nil, false
	}

    discountPriceFloat := 0.0
	if discountPrice != "" {
		_, discountPriceFloat, ok = p.Float(logger, discountPrice)
	}
	if !ok {
		logger.DEBUG_WARN("Failed to convert discount price '%s', will use '%f'", discountPrice, discountPriceFloat)
	}

	discountPricePerUnit := discountPriceFloat
	if discountPriceFloat != 0.0 {
		discountPricePerUnit = ( (discountPriceFloat / priceFloat) * pricePerUnitFloat )
	} 

	product = initProduct(logger, currency, seller, id, name, priceFloat, pricePerUnitFloat, discountPriceFloat, discountPricePerUnit, DiscountPriceInWords, unitType, url, imgURL)

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


func initProduct(logger *logger.Logger, currency, seller, id, name string, price, pricePerUnit, discountPrice, discountPricePerUnit float64, DiscountPriceInWords, UnitType, url, imgURL string) (product *Product) {
	product = new(Product)
	product.Seller = seller
	product.ID = id
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
