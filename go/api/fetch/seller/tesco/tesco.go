package tesco

import (
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/http/url"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func Fetch(logger *logger.Logger, searchValue string) (products *[]*product.Product, ok bool) {

	URL := "https://www.tesco.ie"
	URL += "/groceries/en-IE/search?query="
	URL += searchValue
	URL += "&page=1&count=90"

	waitForJavaScript := false

	urlContext := url.NewUrlContext(URL, waitForJavaScript, fetchFunction)

	logger.INFO("Getting Tesco Products for URL -> %s", URL)

	products, ok = urlContext.Get(logger)
	if !ok {
		logger.ERROR("Failed to get results from Tesco")
		return nil, false
	}

	return products, ok
	
}


func fetchFunction(logger *logger.Logger, doc *goquery.Document, urlContext *url.UrlContext) (products *[]*product.Product, ok bool) {

	productListItems := doc.Find("ul.product-list > li")

	products = &[]*product.Product{}
	productListItems.Each(func(i int, s *goquery.Selection) {

		productName, price, subPrice, specialPrice, productLink, imageURL, ok := parseProductFields(logger, s)
		if !ok {
			logger.DEBUG_WARN("Failed to parse product fields")
			return
		}

		who := "Tesco"
		result, ok := product.NewProduct(logger, who, "ID", productName, price, subPrice, specialPrice, (urlContext.URL + productLink), imageURL)
		if !ok {
			logger.DEBUG_WARN("Failed to create product using name %s, price %s, subPrice %s, specialPrice %s, link %s, imageURL %s", productName, price, subPrice, specialPrice, productLink, imageURL)
			return
		}
		*products = append(*products, result)

	})

	logger.INFO("Tesco - Found %d/%d relevant products", len(*products), productListItems.Length())

	// Return a pointer to the results slice
	return products, true
}



func parseProductFields(logger *logger.Logger, s *goquery.Selection) (name, price, subPrice, specialPrice, link, imageURL string, ok bool) {
	linkSelector := s.Find("a")
	name = s.Find("[data-auto=\"product-tile--title\"]").Text()
	link, exists := linkSelector.Attr("href")

	if !exists {
		logger.DEBUG_WARN("Failed to find link for product")
	}

	price = s.Find(".beans-price__text").Text()
	subPrice = s.Find(".beans-price__subtext").Text()
	specialPrice = getSpecialPrice(logger, s)

	imageURL, exists = s.Find("img").Attr("srcset")
	if !exists {
		logger.DEBUG_WARN("Failed to find image for product")
	}
	if len(imageURL) == 0 {
		logger.DEBUG_WARN("Failed to find image for product")
	}
	imageURL = strings.Split(imageURL, " ")[0]

	return name, price, subPrice, specialPrice, link, imageURL, true
}


func getSpecialPrice(logger *logger.Logger, s *goquery.Selection) (specialPrice string) {
	specialPrice = s.Find(".offer-text").Text()

	pattern := `(\d+\.\d+) Clubcard Price`
    regex := regexp.MustCompile(pattern)

    matches := regex.FindAllStringSubmatch(specialPrice, -1)

    // Extract the matched value from the first match
    if len(matches) > 0 && len(matches[0]) > 1 {
        specialPrice = matches[0][1]
    }

	return specialPrice
}