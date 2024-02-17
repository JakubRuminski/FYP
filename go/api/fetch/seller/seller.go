package seller

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/parse/price_parser"
)


type HTMLParser struct {
	sellerName                       string
	productListItemsPattern          string
	productNamePattern               string

	pricePattern                     string
	pricePerUnitPattern              string
	wasPricePattern                  string
	discountPricePattern             string

	discountPriceInWordsRegexPattern string

	productLinkPattern               string
	productLinkAttribute             string
	imageURLPattern                  string
	imageURLAttribute                string
}

func NewHTMLParser( sellerName,
	                productListItemsPattern,
					productNamePattern,
					pricePattern,
					pricePerUnitPattern,
					wasPricePattern,
					discountPricePattern,
					discountPriceInWordsRegexPattern,
					productLinkPattern,
					productLinkAttribute,
					imageURLPattern,
					imageURLAttribute string,
					) *HTMLParser {
	return &HTMLParser{
		sellerName:                       sellerName,
		productListItemsPattern:          productListItemsPattern,
		productNamePattern:               productNamePattern,
		pricePattern:                     pricePattern,
		pricePerUnitPattern:              pricePerUnitPattern,
		wasPricePattern:                  wasPricePattern,
		discountPricePattern:             discountPricePattern,
		discountPriceInWordsRegexPattern: discountPriceInWordsRegexPattern,

		productLinkPattern:               productLinkPattern,
		productLinkAttribute:             productLinkAttribute,
		
		imageURLPattern:                  imageURLPattern,
		imageURLAttribute:                imageURLAttribute,
	}
}


func (parser *HTMLParser) Parse(logger *logger.Logger, doc *goquery.Document) (products *[]*product.Product, ok bool) {
	productListItems := doc.Find(parser.productListItemsPattern)

	index := -1
	products = &[]*product.Product{}
	productListItems.Each(func(i int, s *goquery.Selection) {
		index++

		productName, ok := parse(index, logger, s, parser.productNamePattern)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse product name", index)
			return 
		}

		currency, price, ok := parseFloat(index, logger, parser.pricePattern, s)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse price", index)
			return
		}

		_, wasPrice, ok := parseFloat(index, logger, parser.wasPricePattern, s)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse was price. Ignoring...", index)
			
		}

		_, discountPrice, ok := parseFloat(index, logger, parser.discountPricePattern, s)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse discount price. Ignoring...", index)
		}

		if wasPrice != 0.0 {
			discountPrice = price
			price = wasPrice
		}

		_, pricePerUnit, pricePerUnitUnitType, ok := parseFloatPerUnit(index, logger, parser.pricePerUnitPattern, s)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse price per unit", index)
			return
		}

		discountPricePerUnit := (discountPrice / price) * pricePerUnit
		logger.DEBUG("%v - discountPricePerUnit: %f", index, discountPricePerUnit)

		discountPriceInWords, ok := parseDiscountPriceInWords(index, logger, s, parser.discountPricePattern, parser.discountPriceInWordsRegexPattern)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse discount price in words. Ignoring...", index)
			
		}

		link, ok := parseByAttribute(index, logger, s, parser.productLinkPattern, parser.productLinkAttribute)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse link", index)
			return
		}

		imageURL, ok := parseByAttribute(index, logger, s, parser.imageURLPattern, parser.imageURLAttribute)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse image URL", index)
			return
		}

		if len(imageURL) == 0 {
			logger.DEBUG_WARN("Failed to find image for product")
		}
		imageURL = strings.Split(imageURL, " ")[0]

		// Create a new Product instance and append it to the products slice
		p, ok := product.NewProduct(
			logger,
			parser.sellerName,
			productName,
			currency,
			price,
			pricePerUnit,
			discountPrice,
			discountPricePerUnit,
			pricePerUnitUnitType,
			discountPriceInWords,
			link,
			imageURL,
		)

		if !ok {
			logger.DEBUG_WARN("%v Failed to create product using name %s, price %s, pricePerUnit %s, discountPrice %s, link %s, imageURL %s", index, productName, price, pricePerUnit, discountPrice, link, imageURL)
			return
		}

		*products = append(*products, p)
	})

	return products, true
}

func parseFloat(index int, logger *logger.Logger, pattern string, s *goquery.Selection) (currency string, price float64, ok bool) {
	priceAsString := s.Find(pattern).Text()
	if priceAsString == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", 0.0, false
	}

	currency, price, ok = price_parser.Float(logger, priceAsString)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert string '%s' to float", index, priceAsString)
		return "", 0.0, false
	}

	return currency, price, true
}

func parseFloatPerUnit(index int, logger *logger.Logger, pattern string, s *goquery.Selection) (currency string, pricePerUnit float64, pricePerUnitUnitType string, ok bool) {
	pricePerUnitAsString := s.Find(pattern).Text()
	if pricePerUnitAsString == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", 0.0, "", false
	}

	currency, pricePerUnit, pricePerUnitUnitType, ok = price_parser.FloatPerUnit(logger, pricePerUnitAsString)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert string '%s' to float", index, pricePerUnitAsString)
		return "", 0.0, "", false
	}

	return currency, pricePerUnit, pricePerUnitUnitType, true
}

func parseDiscountPriceInWords(index int, logger *logger.Logger, s *goquery.Selection, pattern, regexPattern string) (discountPriceInWords string, ok bool) {
	discountPriceString := s.Find(pattern).Text()
	if discountPriceString == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", false
	}

	regex := regexp.MustCompile(regexPattern)
	match := regex.FindString(discountPriceString)

	if match == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with regex pattern '%s'", index, regexPattern)
		return "", true
	}

	return match, true
}

func parse(index int, logger *logger.Logger, s *goquery.Selection, pattern string) (string, bool) {
	result := s.Find(pattern).Text()
	if result == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", false
	}

	return result, true
}

func parseByAttribute(index int, logger *logger.Logger, s *goquery.Selection, pattern, attribute string) (string, bool) {
	result, exists := s.Find(pattern).Attr(attribute)

	if !exists {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s' and attribute '%s'", index, pattern, attribute)
		return "", false
	}

	return result, true
}