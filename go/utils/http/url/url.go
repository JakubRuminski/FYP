package url

import (
	"github.com/PuerkitoBio/goquery"
	
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


type UrlContext struct {
	URL		           string
	WaitForJavaScript  bool
	FetchFunc          func( logger *logger.Logger, doc *goquery.Document, urlContext *UrlContext ) (products *[]*product.Product, ok bool)
}

func NewUrlContext( url string,
					waitForJavaScript bool,
					fetchFunc func( logger *logger.Logger, doc *goquery.Document, urlContext *UrlContext ) (products *[]*product.Product, ok bool) ) (newUrlContext *UrlContext) {
	
	newUrlContext = new(UrlContext)
	newUrlContext.URL               = url
	newUrlContext.WaitForJavaScript = waitForJavaScript
	newUrlContext.FetchFunc         = fetchFunc
	
	return newUrlContext
}


func (urlContext *UrlContext) Get( logger *logger.Logger ) ( products *[]*product.Product, ok bool ) {

	return products, true
}