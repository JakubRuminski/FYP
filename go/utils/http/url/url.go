package url

import (
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"

	"github.com/jakubruminski/FYP/go/api/fetch/seller"
	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)


type UrlContext struct {
	URL		           string
	FullURL            string
	WaitForJavaScript  bool
	FetchFunc          func( logger *logger.Logger, doc *goquery.Document, urlContext *UrlContext, htmlParser *seller.HTMLParser ) (products *[]*product.Product, ok bool)
	htmlParser         *seller.HTMLParser
}

func NewUrlContext( url string,
	                fullURL string,
					waitForJavaScript bool,
					fetchFunc func( logger *logger.Logger, doc *goquery.Document, urlContext *UrlContext, htmlParser *seller.HTMLParser ) (products *[]*product.Product, ok bool),
					htmlParser *seller.HTMLParser ) (newUrlContext *UrlContext) {
	
	newUrlContext = new(UrlContext)
	newUrlContext.URL               = url
	newUrlContext.FullURL           = fullURL
	newUrlContext.WaitForJavaScript = waitForJavaScript
	newUrlContext.FetchFunc         = fetchFunc
	newUrlContext.htmlParser        = htmlParser
	
	return newUrlContext
}


func (urlContext *UrlContext) Get( logger *logger.Logger ) ( products *[]*product.Product, ok bool ) {
	doc, ok := getResponse(logger, urlContext)
	if !ok {
		logger.ERROR("Failed to get results for URL -> %s", urlContext.FullURL)
		return nil, false
	}

	products, ok = urlContext.FetchFunc(logger, doc, urlContext, urlContext.htmlParser)
	if !ok {
		logger.ERROR("Failed to get products from document")
		return nil, false
	}

	return products, true
}

func getResponse(logger *logger.Logger, search *UrlContext) (doc *goquery.Document, ok bool) {
	if search.WaitForJavaScript {
		logger.ERROR("Waiting for JavaScript is not implemented yet")
		return nil, false
	} 

	return getResponseDoNotWaitForJavaScript( logger, search )
}


func getResponseDoNotWaitForJavaScript( logger *logger.Logger, search *UrlContext ) (doc *goquery.Document, ok bool) {
	proxyURL, proxyAPIKey := "PROXY_URL", "PROXY_API_KEY"
	ok = env.GetKeys(logger, &proxyURL, &proxyAPIKey)
	if !ok { return nil, false }

	req, err := http.NewRequest("GET", proxyURL, nil)
	if err != nil {
		logger.ERROR("Error creating HTTP request: %v", err)
		return nil, false
	}

	params := url.Values{
		"api_key":     {proxyAPIKey},
		"url":         {search.FullURL},
	}
	req.URL.RawQuery = params.Encode()

	client := &http.Client{}

	logger.INFO("Sending GET request using %s", search.FullURL)
	resp, err := client.Do(req)
	if err != nil {
		logger.ERROR("Error sending GET request: %v", err)
		return nil, false
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		logger.ERROR("Status code error: %s. Request: %s", resp.Status, search.FullURL)
		logger.ERROR("Response: %v", resp)
		return nil, false
	}

	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		logger.ERROR("Error loading HTTP response body: %v", err)
		return nil, false
	}

	return doc, true
}