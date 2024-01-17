package parse

import "strings"


func Find( stringToSearch string, stringsToLookFor []string ) (string) {
	for _, s := range stringsToLookFor {
		if strings.Contains(stringToSearch, s) {
			return s
		}
	}

	return ""
}


func Strip( stringToStrip string, strip []string ) (stripped string) {
	for _, stripString := range strip {
		stringToStrip = strings.ReplaceAll(stringToStrip, stripString, "")
	}

	return stringToStrip
}