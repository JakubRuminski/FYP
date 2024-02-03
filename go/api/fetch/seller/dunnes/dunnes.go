package dunnes

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/http/url"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func Fetch(logger *logger.Logger, searchValue string) (products *[]*product.Product, ok bool) {
	
	URL     := "https://www.dunnesstoresgrocery.com"
	fullURL := URL + "/sm/delivery/rsid/258/results?q="
	fullURL += searchValue
	fullURL += "&page=1&count=90"

	waitForJavaScript := false

	urlContext := url.NewUrlContext(URL, fullURL, waitForJavaScript, fetchFunction)

	logger.INFO("Getting Tesco Products for URL -> %s", URL)

	products, ok = urlContext.Get(logger)
	if !ok {
		logger.ERROR("Failed to get results from Tesco")
		return nil, false
	}

	return products, ok

}

func fetchFunction(logger *logger.Logger, doc *goquery.Document, urlContext *url.UrlContext) (products *[]*product.Product, ok bool) {
	productListItems := doc.Find(".ColListing--1fk1zey")

	products = &[]*product.Product{}
	productListItems.Each(func(i int, s *goquery.Selection) {

		productName, price, subPrice, specialPrice, specialPriceInWords, productLink, imageURL, ok := parseProductFields(logger, s)
		if !ok {
			logger.DEBUG_WARN("Failed to parse product fields")
			return
		}

		who := "Dunnes"
		result, ok := product.NewProduct(logger, who, "ID", productName, price, subPrice, specialPrice, specialPriceInWords, productLink, imageURL)
		if !ok {
			logger.DEBUG_WARN("Failed to create product using name %s, price %s, subPrice %s, specialPrice %s, link %s, imageURL %s", productName, price, subPrice, specialPrice, specialPriceInWords, (urlContext.URL + productLink), imageURL)
			return
		}
		*products = append(*products, result)

	})

	logger.INFO("Dunnes - Found %d/%d relevant products", len(*products), productListItems.Length())

	return products, true
}

func parseProductFields(logger *logger.Logger, s *goquery.Selection) (name, price, subPrice, specialPrice, specialPriceInWords, link, imageURL string, ok bool) {

	s.Find("[class^='ProductCardTitle--']").Each(func(i int, s *goquery.Selection) {
		name = strings.Replace(s.Text(), "Age restricted item", "", -1)
		name = strings.Replace(name, "Open product description", "", -1)
	})

	link, exists := s.Find("article > a").Attr("href")
	if !exists {
		logger.DEBUG_WARN("Failed to find link for product")
	}

	price = s.Find("[class^='ProductCardPrice--']").Text()
	subPrice = s.Find("[class^='ProductCardPriceInfo--']").Text()
	specialPrice = s.Find(".offer-text").Text()
	specialPrice, specialPriceInWords = getSpecialPrice(logger, specialPrice)

	s.Find("img[class*=Image--]").Each(func(i int, ss *goquery.Selection) {
		imageURL = ss.AttrOr("src", "")
	})

	return name, price, subPrice, specialPrice, specialPriceInWords, link, imageURL, true
}

func getSpecialPrice(logger *logger.Logger, s string) (specialPrice, specialPriceInWords string) {

	// Tries to match the pattern "Buy 2 for €10"
	pattern := `Buy \d+ for €?\d+(\.\d+)?`
    regex := regexp.MustCompile(pattern)

    match := regex.FindString(s)

    if match != "" {
        specialPriceInWords = match
    	return "", specialPriceInWords
    }

	// Tries to match the pattern "€10 Clubcard Price"
	pattern = `SAVE €?\d+(\.\d+)?`
    regex = regexp.MustCompile(pattern)

    matches := regex.FindAllStringSubmatch(s, -1)

    if len(matches) > 0 && len(matches[0]) > 1 {
        specialPriceInWords = matches[0][1]
    }

	return specialPrice, specialPriceInWords
}

