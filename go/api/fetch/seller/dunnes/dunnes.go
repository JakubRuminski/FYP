package dunnes

import (
	"github.com/PuerkitoBio/goquery"

	"github.com/jakubruminski/FYP/go/api/fetch/seller"
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

	htmlParser := seller.NewHTMLParser(
		"Dunnes",
		".ColListing--1fk1zey",
		"[class^='ProductCardTitle--']",
		"[class^='ProductCardPrice--']",
		"[class^='ProductCardPriceInfo--']",
		"[class^='WasPrice--']",
		`[data-testid="promotionBadgeComponent-testId"]`,
		`Buy \d+ for €?\d+(\.\d+)?`,
		[]string{},
		"article > a",
		"href",
		"img[class*=Image--]",
		"src",
	)

	urlContext := url.NewUrlContext(URL, fullURL, waitForJavaScript, fetchFunction, htmlParser)

	logger.INFO("Getting Tesco Products for URL -> %s", URL)

	products, ok = urlContext.Get(logger)
	if !ok {
		logger.ERROR("Failed to get results from Tesco")
		return nil, false
	}

	return products, ok

}

func fetchFunction(logger *logger.Logger, doc *goquery.Document, urlContext *url.UrlContext, htmlParser *seller.HTMLParser) (products *[]*product.Product, ok bool) {
	products, ok = htmlParser.Parse(logger, doc)
	if !ok {
		logger.ERROR("Failed to parse products")
		return nil, false
	}

	return products, true
}