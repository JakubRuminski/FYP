package url

import (
	"math/rand"
	"net/http"
	"time"

	"github.com/PuerkitoBio/goquery"

	"github.com/jakubruminski/FYP/go/api/fetch/seller"
	"github.com/jakubruminski/FYP/go/api/product"
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


func getResponseDoNotWaitForJavaScript(logger *logger.Logger, search *UrlContext) (doc *goquery.Document, ok bool) {
	attempts := 2
	// Initialize a http.Client with a cookie jar to maintain session across requests
	client := &http.Client{}

	for i := 0; i < attempts; i++ {
		req, err := http.NewRequest("GET", search.FullURL, nil)
		if err != nil {
			logger.ERROR("Error creating HTTP request: %v", err)
			return nil, false
		}

		// Mimic a browser more closely with updated headers
		req.Header.Set("User-Agent", getRandomUserAgent(logger))
		req.Header.Set("Accept-Language", "en-US,en;q=0.5")
		req.Header.Set("Referer", "https://www.google.com/") // Pretend the request comes from a Google search
		req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")

		// More sophisticated delay to mimic human browsing patterns
		delay := time.Duration(rand.Intn(5)+2) * time.Second // Random delay between 2 and 6 seconds
		time.Sleep(delay)

		logger.INFO("Sending GET request to %s", search.FullURL)

		resp, err := client.Do(req)
		if err != nil {
			logger.ERROR("Error sending GET request: %v", err)
			return nil, false
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				logger.ERROR("Error loading HTTP response body: %v", err)
				return nil, false
			}

			return doc, true
		} else {
			logger.ERROR("Received HTTP status code %d for URL %s", resp.StatusCode, search.FullURL)
		}

		logger.DEBUG("Retrying request to %s", search.FullURL)
	}

	logger.ERROR("Failed to get response after %d attempts", attempts)
	return nil, false
}

func getRandomUserAgent(logger *logger.Logger) string {
	// Create a new source with a seed based on the current time
	source := rand.NewSource(time.Now().UnixNano())

	// Create a new random number generator using the source
	rng := rand.New(source)

	// Updated and expanded list of User-Agent strings, including both desktop and mobile browsers
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.1.2 Safari/605.1.15",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:97.0) Gecko/20100101 Firefox/97.0",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 Edg/98.0.1108.56",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36",
		"Mozilla/5.0 (iPad; CPU OS 13_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 13_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/13.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Linux; Android 10; SM-A505FN) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Mobile Safari/537.36",
		"Mozilla/5.0 (Linux; Android 10; SM-G981B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Mobile Safari/537.36 EdgA/45.12.4.5121",
		"Mozilla/5.0 (Linux; U; Android 10; en-us; Redmi Note 8 Pro) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/78.0.3904.108 UCBrowser/13.2.5.1300 Mobile Safari/537.36",
	}

	// Get a random index from the list
	randomIndex := rng.Intn(len(userAgents))

	// Return the randomly selected User-Agent string
	userAgent := userAgents[randomIndex]
	logger.DEBUG("Random User-Agent: %s", userAgent)
	return userAgent
}