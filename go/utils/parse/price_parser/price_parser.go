package price_parser

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/parse"
	"github.com/jakubruminski/FYP/go/utils/slice"
)


func Float(logger *logger.Logger, price string) (currency string, priceFloat float64, ok bool) {
	currency, price, ok = findAndStripCurrency(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to find and strip currency from '%s'", price)
		return "", 0.0, false
	}

	price = parse.StripNonNumeric(price)

	priceFloat, ok = convertToFloat(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to convert string '%s' to float", price)
		return "", 0.0, false
	}
	
	return currency, priceFloat, true
}


func FloatPerUnit(logger *logger.Logger, price string) (currency string, priceFloat float64, pricePerUnit string, ok bool) {
	currency, price, ok = findAndStripCurrency(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to find and strip currency from '%s'", price)
		return "", 0.0, "", false
	}

	parsedPrice, perUnitQuantity, perUnit, measurement, ok := stripMeasurement(logger, price)
	if !ok {
		logger.DEBUG_WARN("Failed to strip price '%s'", price)
		return currency, 0.0, "", false
	}
	logger.INFO("parsedPrice: %s, perUnitQuantity: %s, perUnit: %s, measurement: %s", parsedPrice, perUnitQuantity, perUnit, measurement)
	priceFloat, ok = convertPrice(logger, parsedPrice, perUnitQuantity, perUnit)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price '%s'", price)
		return currency, 0.0, "", false
	}
	
	
	return currency, priceFloat, measurement, true
}


func findAndStripCurrency(logger *logger.Logger, price string) (currency, priceStripped string, ok bool) {
	currency = parse.Find(price, []string{"€", "£", "$"})
	price   = parse.Strip(price, []string{"€", "£", "$"})
	return currency, price, true
}


func convertToFloat(logger *logger.Logger, price string) (priceFloat float64, ok bool) {
	priceFloat, err := strconv.ParseFloat(price, 64)
	if err != nil {
		logger.DEBUG_WARN("Error converting string '%s' to float", price)
		return 0.0, false
	}
	
	return priceFloat, true
}


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


func stripMeasurement(logger *logger.Logger, price string) (parsedPrice, perUnitQuantity, perUnit, measurement string, ok bool) {
	var parsedPriceArray []string  // will look like this later -> ["700", "70", "cl"]
	if strings.Contains(price, "/") {
		parsedPriceArray = strings.Split(price, "/")
	} else if strings.Contains(price, "per") {
		parsedPriceArray = strings.Split(price, "per")
	}

	if len(parsedPriceArray) != 2 {
		logger.DEBUG_WARN("Failed to parse price '%s'.", price)
		return "", "", "", "", false
	}

	regexNumeric := regexp.MustCompile(`[0-9.]+`)
	regexNominal := regexp.MustCompile(`[^0-9.\s]+`)

	for measurement, measurementDict := range unitTypes_DICT {
		for _, unitTypeSlice := range measurementDict {
			for _, unit := range unitTypeSlice {
				parsedPrice = regexNumeric.FindString(parsedPriceArray[0])
				perUnitQuantity = regexNumeric.FindString(parsedPriceArray[1])
				perUnit = regexNominal.FindString(parsedPriceArray[1])

				if parsedPrice == "" || perUnit == "" {
					logger.DEBUG("parsed Price '%s' or perUnit '%s' was empty. Skipping...", parsedPrice, perUnit)
					continue
				}
				if !strings.Contains(perUnit, unit) {
					continue
				}
		
				if perUnitQuantity == "" {
					perUnitQuantity = "1"
				}

				perUnit = unit
				return parsedPrice, perUnitQuantity, perUnit, measurement, true
			}
		}
	}
	
	logger.DEBUG("perUnit '%s' did not contain any recognised Unit Type. Skipping...", perUnit)
	return "", "", "unknown", "unknown", false
}


func convertPrice(logger *logger.Logger, parsedPrice, perUnitQuantity, unitType string) (priceFloat float64, ok bool) {
	parsedPriceFloat, ok := convertToFloat(logger, parsedPrice)
	if !ok {
		logger.DEBUG_WARN("Failed to convert price per unit '%s' to float", parsedPrice)
		return 0.0, false
	}

	parsedUnitQuantityFloat, ok := convertToFloat(logger, perUnitQuantity)
	if !ok {
		logger.DEBUG_WARN("Failed to convert unit quantity '%s' to float", perUnitQuantity)
		return 0.0, false
	}

	isKiloOrLitre := slice.ContainsString(unitTypes_DICT["kilogram"]["kilogram"], unitType) || slice.ContainsString(unitTypes_DICT["litre"]["litre"], unitType)
	isGramOrMillilitre := slice.ContainsString(unitTypes_DICT["kilogram"]["gram"], unitType) || slice.ContainsString(unitTypes_DICT["litre"]["millilitre"], unitType)
	isCentilitre := slice.ContainsString(unitTypes_DICT["litre"]["centilitre"], unitType)
	isEach := slice.ContainsString(unitTypes_DICT["each"]["each"], unitType)

	if isKiloOrLitre {
		return parsedPriceFloat, true

	} else if isGramOrMillilitre {
		return (parsedPriceFloat * (1000 / parsedUnitQuantityFloat)), true

	} else if isCentilitre {
		return (parsedPriceFloat * (100 / parsedUnitQuantityFloat)), true

	} else if isEach {
		return (parsedPriceFloat * parsedUnitQuantityFloat), true
	}


	logger.DEBUG_WARN("Failed to convert price per unit type '%s'", parsedPrice)
	return 0.0, false
}

