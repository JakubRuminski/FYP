package price_parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/parse"
	"github.com/jakubruminski/FYP/go/utils/slice"
)


func Float(index int, logger *logger.Logger, price string) (currency string, priceFloat float64, ok bool) {
	currency, price, ok = findAndStripCurrency(index, logger, price)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to find and strip currency from '%s'", index, price)
		return "", 0.0, false
	}

	price = strings.Replace(price, " ", "", -1)

	priceFloat, ok = convertToFloat(index, logger, price)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert string '%s' to float", index, price)
		return "", 0.0, false
	}
	
	return currency, priceFloat, true
}


func FloatPerUnit(index int, logger *logger.Logger, price string) (currency string, priceFloat float64, pricePerUnit string, ok bool) {
	currency, price, ok = findAndStripCurrency(index, logger, price)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to find and strip currency from '%s'", index, price)
		return "", 0.0, "", false
	}

	parsedPrice, perUnitQuantity, perUnit, measurement, ok := stripMeasurement(index, logger, price)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to strip price '%s'", index, price)
		return currency, 0.0, "", false
	}
	priceFloat, ok = convertPrice(index, logger, parsedPrice, perUnitQuantity, perUnit)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert price '%s'", index, price)
		return currency, 0.0, "", false
	}	
	
	return currency, priceFloat, measurement, true
}


func findAndStripCurrency(index int, logger *logger.Logger, price string) (currency, priceStripped string, ok bool) {
	currency = parse.Find(price, []string{"€", "£", "$"})
	price   = parse.Strip(price, []string{"€", "£", "$"})
	return currency, price, true
}


func convertToFloat(index int, logger *logger.Logger, price string) (priceFloat float64, ok bool) {
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logger.DEBUG_WARN("%v - Error converting string '%s' to float", index, price)
		return 0.0, false
	}
	
	return priceFloat, true
}

// Good To Know: The order of the unit types in the map is ignored when iterated over.
// Sometimescausing kilograms to be stepped over first of grams, and vice versa.
//
var unitTypes_DICT = map[string]map[string][]string{
    "kilogram": {
        "kilogram": {"kilograms", "kilogram", "kilo", "kg"},
        "gram":     {"grams", "gram", "g"},
    },
    "litre": {
        "centilitre": {"centilitres", "centilitre", "cl"},
        "millilitre": {"millilitres", "millilitre", "ml"},
        "litre":      {"litres", "litre", "l"},
    },
    "each": {
        "each": {"each", "unit", "item", "items", "sht"},
    },
}


func stripMeasurement(index int, logger *logger.Logger, price string) (parsedPrice, perUnitQuantity, perUnit, measurement string, ok bool) {
	var parsedPriceArray []string  // will look like this later -> ["700", "70", "cl"]
	if strings.Contains(price, "/") {
		parsedPriceArray = strings.Split(price, "/")
	} else if strings.Contains(price, "per") {
		parsedPriceArray = strings.Split(price, "per")
	} else {
		parsedPriceArray = strings.Split(price, " ")
		if len(parsedPriceArray) == 2 && !slice.ContainsString(unitTypes_DICT["each"]["each"], parsedPriceArray[1]) {
            logger.DEBUG_WARN("%v - Failed to split price '%s' by '/' or 'per'.", index, price)
			return "", "", "", "", false
		}
	}

	if len(parsedPriceArray) != 2 {
		logger.DEBUG_WARN("%v - Failed to parse price '%s'.", index, price)
		return "", "", "", "", false
	}

	regexNumeric := regexp.MustCompile(`[0-9.]+`)
	regexNominal := regexp.MustCompile(`[^0-9.\s]+`)

	logger.DEBUG("%v - Analysing string array '%v'", index, parsedPriceArray)

	for measurement, measurementDict := range unitTypes_DICT {
		for _, unitTypeSlice := range measurementDict {
			for _, unit := range unitTypeSlice {
				parsedPrice = regexNumeric.FindString(parsedPriceArray[0])
				perUnitQuantity = regexNumeric.FindString(parsedPriceArray[1])
				perUnit = regexNominal.FindString(parsedPriceArray[1])

				if parsedPrice == "" || perUnit == "" {
					logger.DEBUG("%v - parsed Price '%s' or perUnit '%s' was empty. Skipping...", index, parsedPrice, perUnit)
					continue
				}
				if !strings.Contains(perUnit, unit) {
					continue
				}
		
				if perUnitQuantity == "" {
					perUnitQuantity = "1"
				}

				logger.DEBUG("%v - Found parsed Price '%s'", index, parsedPrice)
				logger.DEBUG("%v - Found perUnitQuantity '%s'", index, perUnitQuantity)
				logger.DEBUG("%v - Found perUnit '%s'", index, perUnit)
				logger.DEBUG("%v - Found measurement '%s'", index, measurement)

				// Do not uncomment this line, it causes instable unit type recognition
				// perUnit = unit

				return parsedPrice, perUnitQuantity, perUnit, measurement, true
			}
		}
	}
	
	logger.DEBUG("%v - perUnit '%s' did not contain any recognised Unit Type. Skipping...", index, perUnit)
	return "", "", "unknown", "unknown", false
}


func convertPrice(index int, logger *logger.Logger, parsedPrice, perUnitQuantity, unitType string) (priceFloat float64, ok bool) {
	parsedPriceFloat, ok := convertToFloat(index, logger, parsedPrice)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert price per unit '%s' to float", index, parsedPrice)
		return 0.0, false
	}

	parsedUnitQuantityFloat, ok := convertToFloat(index, logger, perUnitQuantity)
	if !ok {
		logger.DEBUG_WARN("%v - Failed to convert unit quantity '%s' to float", index, perUnitQuantity)
		return 0.0, false
	}

	isCentilitre := (slice.ContainsString(unitTypes_DICT["litre"]["centilitre"], unitType))
	logger.DEBUG("%v - isCentilitre: %v", index, isCentilitre)

	isGramOrMillilitre := (slice.ContainsString(unitTypes_DICT["kilogram"]["gram"], unitType) || slice.ContainsString(unitTypes_DICT["litre"]["millilitre"], unitType))
	logger.DEBUG("%v - isGramOrMillilitre: %v", index, isGramOrMillilitre)

	isKiloOrLitre := (slice.ContainsString(unitTypes_DICT["kilogram"]["kilogram"], unitType) || slice.ContainsString(unitTypes_DICT["litre"]["litre"], unitType))
	logger.DEBUG("%v - isKiloOrLitre: %v", index, isKiloOrLitre)
	
	isEach := slice.ContainsString(unitTypes_DICT["each"]["each"], unitType)
	logger.DEBUG("%v - isEach: %v", index, isEach)

	
	if isCentilitre {
		return (parsedPriceFloat * (100 / parsedUnitQuantityFloat)), true

	} else if isGramOrMillilitre {
		return (parsedPriceFloat * (1000 / parsedUnitQuantityFloat)), true

	} else if isKiloOrLitre {
		return parsedPriceFloat, true

	} else if isEach {
		return (parsedPriceFloat * parsedUnitQuantityFloat), true
	}


	logger.DEBUG_WARN("%v - Failed to convert price per unit type '%s'", index, parsedPrice)
	return 0.0, false
}

