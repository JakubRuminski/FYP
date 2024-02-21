package supervalu

import (
	"github.com/PuerkitoBio/goquery"
	
	"github.com/jakubruminski/FYP/go/api/fetch/seller"
	"github.com/jakubruminski/FYP/go/api/product"
	
	"github.com/jakubruminski/FYP/go/utils/http/url"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


func Fetch(logger *logger.Logger, searchValue string) (products *[]*product.Product, ok bool) {

	URL     := "https://shop.supervalu.ie"
	fullURL := URL + "/sm/delivery/rsid/5550/results?q="
	fullURL += searchValue
	fullURL += "" // TODO: How can this work: "&page=1&count=90"

	waitForJavaScript := false

	htmlParser := seller.NewHTMLParser(
		"SuperValu",
		
		"[class^='ColListing--']",  // DONE
		
		"[class^='ProductCardTitle--']",
		[]string{"Open product description"},  // Product Name Strings to Strip
		
		"[class^='ProductCardPrice--']",
		"[class^='ProductCardPriceInfo--']",
		"[class^='WasPrice--']",
		[]string{"was"},
		
		`[data-testid="promotionBadgeComponent-testId"]`,
		"",
		[]string{`\d+ for €?\d+(\.\d+)?`},
		[]string{"ONLY", "SAVE"},
		
		`\d+ for €?\d+(\.\d+)?`,
		[]string{},
		
		"",
		"article > a",
		"href",
		
		"img[class^=ProductCardImage--]",
		"src",
	)

	urlContext := url.NewUrlContext(URL, fullURL, waitForJavaScript, fetchFunction, htmlParser)

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
