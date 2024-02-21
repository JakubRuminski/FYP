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
	sellerName                           string

	productListItemsPattern              string
	productNamePattern                   string
	productNameStringsToStrip            []string

	pricePattern                         string
	pricePerUnitPattern                  string
	wasPricePattern                      string
	wasPriceStringToStrip                []string

	discountPricePattern                 string
	discountPriceAllowedRegexPattern     string
	discountPriceProhibitedRegexPattern  []string
	discountPriceStringToStrip           []string

	discountPriceInWordsRegexPattern     string
	discountPriceInWordsStringsToStrip   []string

	productLinkPrefix					 string
	productLinkPattern                   string
	productLinkAttribute                 string
	imageURLPattern                      string
	imageURLAttribute                    string
}

func NewHTMLParser( sellerName,
	                productListItemsPattern,
					productNamePattern string,
					productNameStringsToStrip []string,
					pricePattern,
					pricePerUnitPattern,
					wasPricePattern string,
					wasPriceStringToStrip []string,
					discountPricePattern string,
					discountPriceAllowedRegexPattern string,
					discountPriceProhibitedRegexPattern []string,
					discountPriceStringToStrip []string,
					discountPriceInWordsRegexPattern string,
					discountPriceInWordsStringsToStrip []string,
					productLinkPrefix,
					productLinkPattern,
					productLinkAttribute,
					imageURLPattern,
					imageURLAttribute string,
					) *HTMLParser {
	return &HTMLParser{
		sellerName:                           sellerName,
		productListItemsPattern:              productListItemsPattern,
		
		productNamePattern:                   productNamePattern,
		productNameStringsToStrip:            productNameStringsToStrip,
		
		pricePattern:                         pricePattern,
		pricePerUnitPattern:                  pricePerUnitPattern,
		
		wasPricePattern:                      wasPricePattern,
		wasPriceStringToStrip:                wasPriceStringToStrip,

		discountPricePattern:                 discountPricePattern,
		discountPriceAllowedRegexPattern:     discountPriceAllowedRegexPattern,
		discountPriceProhibitedRegexPattern:  discountPriceProhibitedRegexPattern,
		discountPriceStringToStrip:           discountPriceStringToStrip,

		discountPriceInWordsRegexPattern:     discountPriceInWordsRegexPattern,
		discountPriceInWordsStringsToStrip:   discountPriceInWordsStringsToStrip,

		productLinkPrefix:					  productLinkPrefix,
		productLinkPattern:                   productLinkPattern,
		productLinkAttribute:                 productLinkAttribute,
		
		imageURLPattern:                      imageURLPattern,
		imageURLAttribute:                    imageURLAttribute,
	}
}


func (parser *HTMLParser) Parse(logger *logger.Logger, doc *goquery.Document) (products *[]*product.Product, ok bool) {
	productListItems := doc.Find(parser.productListItemsPattern)

	index := -1
	products = &[]*product.Product{}
	productListItems.Each(func(i int, s *goquery.Selection) {
		index++
		logger.DEBUG("%v - START -----------------", index)
		defer logger.DEBUG("%v - END -----------------", index)

		// print statement which displays the raw html elements and inner html of the product
		rawHTML, err := s.Html()
		if err != nil {
			logger.WARN("%v - Failed to get raw html", index)
			return
		}
		
		logger.DEBUG("%v - RAW HTML:\n%s", index, rawHTML)
		

		productName, ok := parse(index, logger, s, parser.productNamePattern, parser.productNameStringsToStrip...)
		if !ok {
			logger.WARN("%v - Failed to parse product name", index)
			return 
		}

		link, ok := parseByAttribute(index, logger, s, parser.productLinkPattern, parser.productLinkAttribute)
		if !ok {
			logger.WARN("%v - Failed to parse link", index)
			return
		}

		if parser.productLinkPrefix != "" {
			link = parser.productLinkPrefix + link
		}

		currency, price, ok := parseFloat(index, logger, parser.pricePattern, s, false, "", []string{}, []string{})
		if !ok {
			logger.WARN("%v - [%s] Failed to parse price for product", index, link)
			return
		}

		_, wasPrice, ok := parseFloat(index, logger, parser.wasPricePattern, s, true, "", []string{}, parser.wasPriceStringToStrip)
		if !ok {
			logger.DEBUG_WARN("%v - [%s] Failed to parse was price. Ignoring...", index, link)
			
		}

		_, discountPrice, ok := parseFloat(index, logger, parser.discountPricePattern, s, true, parser.discountPriceAllowedRegexPattern, parser.discountPriceProhibitedRegexPattern, parser.discountPriceStringToStrip)
		if !ok {
			logger.DEBUG_WARN("%v - Failed to parse discount price. Ignoring... [%s] ", index, link)
		}

		if wasPrice != 0.0 {
			discountPrice = price
			price = wasPrice
		}

		_, pricePerUnit, pricePerUnitUnitType, ok := parseFloatPerUnit(index, logger, parser.pricePerUnitPattern, s)
		if !ok {
			logger.WARN("%v - [%s] Failed to parse price per unit", index, link)
			return
		}

		discountPricePerUnit := (discountPrice / price) * pricePerUnit

		discountPriceInWords, ok := parseDiscountPriceInWords(index, logger, s, parser.discountPricePattern, parser.discountPriceInWordsRegexPattern, parser.discountPriceInWordsStringsToStrip...)
		if !ok {
			logger.DEBUG_WARN("%v - [%s] Failed to parse discount price in words. Ignoring...", index, link)
		}

		imageURL, ok := parseByAttribute(index, logger, s, parser.imageURLPattern, parser.imageURLAttribute)
		if !ok {
			logger.WARN("%v - [%s] Failed to parse image URL", index, link)
			return
		}

		if len(imageURL) == 0 {
			logger.WARN("%v - [%s] Failed to find image for product", index, link)
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

		logger.DEBUG("%v - Parsed p.Seller: %s", index, p.Seller)
		logger.DEBUG("%v - Parsed p.Name: %s", index, p.Name)
		logger.DEBUG("%v - Parsed p.Currency: %s", index, p.Currency)
		logger.DEBUG("%v - Parsed p.Price: %f", index, p.Price)
		logger.DEBUG("%v - Parsed p.PricePerUnit: %f", index, p.PricePerUnit)
		logger.DEBUG("%v - Parsed p.DiscountPrice: %f", index, p.DiscountPrice)
		logger.DEBUG("%v - Parsed p.DiscountPricePerUnit: %f", index, p.DiscountPricePerUnit)
		logger.DEBUG("%v - Parsed p.DiscountPriceInWords: %s", index, p.DiscountPriceInWords)
		logger.DEBUG("%v - Parsed p.URL: %s", index, p.URL)
		logger.DEBUG("%v - Parsed p.ImgURL: %s", index, p.ImgURL)

		// logger.DATA(`unOrderedProducts = append(unOrderedProducts, &Product{ Seller: "%s", Name: "%s", Currency: "%s", Price: %f, PricePerUnit: %f, DiscountPrice: %f, DiscountPricePerUnit: %f, DiscountPriceInWords: "%s", URL: "%s", ImgURL: "%s" })`, p.Seller, p.Name, p.Currency, p.Price, p.PricePerUnit, p.DiscountPrice, p.DiscountPricePerUnit, p.DiscountPriceInWords, p.URL, p.ImgURL) 

		*products = append(*products, p)
	})

	logger.INFO("Successfully parsed %v out of %v products", len(*products), index+1)

	return products, true
}

func parseFloat(index int, logger *logger.Logger, pattern string, s *goquery.Selection, optional bool, allowedRegexPattern string, prohibitedRegexPattern, stringsToStrip []string) (currency string, price float64, ok bool) {
	priceAsString := s.Find(pattern).Text()
	if priceAsString == "" && optional {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'. The optional flag was set to %v", index, pattern, optional)
		return "", 0.0, true

	} else if priceAsString == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", 0.0, false
	}

	for _, p := range prohibitedRegexPattern {
		regex := regexp.MustCompile(p)
		matches := regex.FindAllStringSubmatch(priceAsString, -1)
		if len(matches) > 0 {
			logger.DEBUG_WARN("%v - Found matches with prohibited regex pattern '%s'", index, p)
			return "", 0.0, false
		}
	}

	if allowedRegexPattern != "" {
		regex := regexp.MustCompile(allowedRegexPattern)
		matches := regex.FindAllStringSubmatch(priceAsString, -1)
		if len(matches) == 0 {
			logger.DEBUG_WARN("%v - Failed to find anything with regex pattern '%s'", index, allowedRegexPattern)
			return "", 0.0, false
		}

		if len(matches) > 0 && len(matches[0]) > 1 {
			logger.DEBUG("%v - Found matches: %v", index, matches[0])
			priceAsString = matches[0][1]
		}
	}

	priceAsString = strings.ToLower(priceAsString)
	for _, s := range stringsToStrip {
		s = strings.ToLower(s)
		if strings.Contains(priceAsString, s) {
			priceAsString = strings.Replace(priceAsString, s, "", -1)
			logger.DEBUG("%v - Stripped %s from %s", index, s, priceAsString)
		}
	}

	currency, price, ok = price_parser.Float(index, logger, priceAsString)
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

	currency, pricePerUnit, pricePerUnitUnitType, ok = price_parser.FloatPerUnit(index, logger, pricePerUnitAsString)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert string '%s' to float", index, pricePerUnitAsString)
		return "", 0.0, "", false
	}

	return currency, pricePerUnit, pricePerUnitUnitType, true
}

func parseDiscountPriceInWords(index int, logger *logger.Logger, s *goquery.Selection, pattern, regexPattern string, stringsToStrip ...string) (string, bool) {
	discountPriceString := s.Find(pattern).Text()
	if discountPriceString == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", true
	}

	regex := regexp.MustCompile(regexPattern)
	matches := regex.FindAllStringSubmatch(discountPriceString, -1)
	if len(matches) == 0 {
		logger.DEBUG_WARN("%v - Failed to find anything with regex pattern '%s'", index, regexPattern)
		return "", true
	}

	if len(matches) > 0 && len(matches[0]) > 1 {
		discountPriceString = matches[0][0]
		logger.DEBUG("%v - Found matches: %v", index, discountPriceString)
	}

	for _, s := range stringsToStrip {
		discountPriceString = strings.ToLower(discountPriceString)
		s = strings.ToLower(s)

		if strings.Contains(discountPriceString, s) {
			discountPriceString = strings.Replace(discountPriceString, s, "", -1)
			logger.DEBUG("%v - Stripped '%s' from '%s'", index, s, discountPriceString)
		}
	}

	return discountPriceString, true
}

func parse(index int, logger *logger.Logger, s *goquery.Selection, pattern string, stringsToStrip ...string) (string, bool) {
	result := s.Find(pattern).Text()
	if result == "" {
		logger.DEBUG_WARN("%v - Failed to find anything with pattern '%s'", index, pattern)
		return "", false
	}

	for _, s := range stringsToStrip {
		result = strings.ToLower(result)
		s = strings.ToLower(s)
		if strings.Contains(result, s) {
			result = strings.Replace(result, s, "", -1)
			logger.DEBUG("%v - Stripped %s from %s", index, s, result)
		}
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