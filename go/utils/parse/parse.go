package parse

import (
	"regexp"
	"strings"
)


func Find( stringToSearch string, stringsToLookFor []string ) (string) {
	for _, s := range stringsToLookFor {
		if strings.Contains(stringToSearch, s) {
			return s
		}
	}

	return ""
}

// StripNonNumeric removes all non-numeric characters from a string
// 
// Exceptions:
//   - decimal points
//
func StripNonNumeric( stringToStrip string ) (stripped string) {
	reg := regexp.MustCompile("[^0-9.]")

	return reg.ReplaceAllString(stringToStrip, "")
}


func Strip( stringToStrip string, strip []string ) (stripped string) {
	for _, stripString := range strip {
		stringToStrip = strings.ReplaceAll(stringToStrip, stripString, "")
	}

	return stringToStrip
}