package aldi

// import (
// 	"encoding/json"
// 	"strconv"
// 	"strings"

// 	"github.com/PuerkitoBio/goquery"
// 	"github.com/jakubruminski/FYP/go/api/product"
// 	"github.com/jakubruminski/FYP/go/utils/http/url"
// 	"github.com/jakubruminski/FYP/go/utils/logger"
// )


// func SearchAldi(logger *logger.Logger, searchValue string) (products *[]*product.Product, ok bool) {

// 	URL     := "https://groceries.aldi.ie/en-GB"
// 	fullURL := URL + "/Search?keywords="
// 	fullURL += searchValue
// 	fullURL += "&page=1"// Aldi doesn't have pagination

// 	waitForJavaScript := false

// 	urlContext := url.NewUrlContext(URL, fullURL, waitForJavaScript, fetchFunction)

// 	logger.INFO("Getting Tesco Products for URL -> %s", URL)

// 	products, ok = urlContext.Get(logger)
// 	if !ok {
// 		logger.ERROR("Failed to get results from Tesco")
// 		return nil, false
// 	}

// 	return products, ok
// }

// func fetchFunction(logger *logger.Logger, doc *goquery.Document, urlContext *url.UrlContext) (products *[]*product.Product, ok bool) {

// 	productListItems, exists := doc.Find("#searchResults").Attr("data-context")
// 	if !exists { logger.DEBUG_WARN("Failed to find link for product") }

// 	var products productsWrapper
// 	err := json.Unmarshal([]byte(productListItems), &products)
// 	if err != nil {
// 		logger.ERROR("Error occurred during unmarshalling. Err: %v", err)
// 		return nil, false
// 	}

// 	logger.INFO("Aldi - Found %d products", len(products.AldiProducts))

// 	res := []*r.Result{}
// 	for _, item := range products.AldiProducts {
// 		if rate, ok := currency.Rates["UK"]["rate"].(float64); ok {
// 			item.ListPrice = item.ListPrice / rate
// 		} else {
// 			logger.ERROR("Error: rate is not a float64")
// 		}
// 		price := strconv.FormatFloat(item.ListPrice, 'f', 2, 64)
// 		item.SizeVolume = price + " per " + item.SizeVolume

// 		item.Url = "https://groceries.aldi.ie/en-GB" + item.Url

// 		logger.DEBUG("Retrieved productName: %s", item.FullDisplayName)
//         logger.DEBUG("Retrieved productLink: %s", (item.Url))
//         logger.DEBUG("Retrieved price: %s", price)
//         logger.DEBUG("Retrieved subPrice: %s", item.SizeVolume)
//         logger.DEBUG("Retrieved specialPrice: %s", (""))
//         logger.DEBUG("Retrieved imageURL: %s", (item.ImageUrl))

// 		who := "Aldi"
// 		result := r.NewResult(who, item.FullDisplayName, price, item.SizeVolume, "", item.Url, item.ImageUrl)
// 		res = append(res, result)
// 	}

// 	return &res, true
// }

// // -------- STRUCTS ------------------------------------------------------
// type aldiProduct struct {
// 	FullDisplayName string  `json:"FullDisplayName"`
//     ListPrice       float64 `json:"ListPrice"`
//     Url             string  `json:"Url"`
//     ImageUrl        string  `json:"ImageUrl"`
//     SizeVolume      string  `json:"SizeVolume"`
// }

// type productsWrapper struct {
//     AldiProducts []aldiProduct `json:"SearchResults"`
// }
// // -----------------------------------------------------------------------


// // -------- UTILS --------------------------------------------------------
// func multipleWords(logger *logger.Logger, searchValue string) string {
// 	newSearchValue := searchValue
// 	if strings.Contains(searchValue, " ") {
// 		newSearchValue = strings.Replace(searchValue, " ", "+", -1)
// 	}
// 	logger.DEBUG("SearchValue '%s' - NewSearchValue '%s'", searchValue, newSearchValue)

// 	return newSearchValue
// }

// // -----------------------------------------------------------------------