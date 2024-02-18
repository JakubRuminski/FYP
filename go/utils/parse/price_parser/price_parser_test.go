package price_parser

import (
	"testing"

	"github.com/jakubruminski/FYP/go/utils/logger"
)

func TestFloat(t *testing.T) {
	logger := &logger.Logger{}

	testCases := []struct {
		priceAsString             string
		expectedCurrency          string
		expectedPriceType         string
		expectedPrice             float64
	}{
		{"€43 per 70cl",  "€", "litre", 61.42857142857143},
		{"€50 per 100cl", "€", "litre", 50.0},
		{"€5/kg",         "€", "kilogram", 5.0},
		{"€2/g",          "€", "kilogram", 2000.0},
		{"€9.68/l",       "€", "litre", 9.68},
		
		{"€3/litre",      "€", "litre", 3.0},
		{"€3/ml",         "€", "litre", 3000.0},
		{"€0.01/cl",      "€", "litre", 1.0},
		
		{"€2/item",       "€", "each", 2.0},

		// These should fail
		{"€5",            "€", "", 0.0},
		{"€5.00",         "€", "", 0.0},

	}

	success := true
	for index, tc := range testCases {
		index++

		currency, priceFloat, pricePerUnit, ok := FloatPerUnit(index, logger, tc.priceAsString)
		if tc.expectedPrice == 0.0 && ok {
			t.Errorf("Expected ok to be false, but got true for test case PriceType: %s, Price: %s", tc.expectedPriceType, tc.priceAsString)
			success = false
		} else if tc.expectedPrice != 0.0 && !ok {
			t.Errorf("Expected ok to be true, but got false for test case PriceType: %s, Price: %s", tc.expectedPriceType, tc.priceAsString)
			success = false
		}

		if currency != tc.expectedCurrency {
			t.Errorf("Expected %s, but got %s for test case PriceType: %s, Price: %s", tc.expectedCurrency, currency, tc.expectedPriceType, tc.priceAsString)
			success = false
		}
		
		if pricePerUnit != tc.expectedPriceType {
			t.Errorf("Expected %s, but got %s for test case PriceType: %s, Price: %s", tc.expectedPriceType, pricePerUnit, tc.expectedPriceType, tc.priceAsString)
			success = false
		}

		if priceFloat != tc.expectedPrice {
			t.Errorf("Expected %f, but got %f for test case PriceType: %s, Price: %s", tc.expectedPrice, priceFloat, tc.expectedPriceType, tc.priceAsString)
			success = false
		}

		logger.INFO("%d FINISHED - PriceType: %s, Price: %s, Result: %f, Expected: %f", index, tc.expectedPriceType, tc.priceAsString, priceFloat, tc.expectedPrice)
	}

	if success {
		logger.INFO("All %d tests passed", len(testCases))
	} else {
		logger.ERROR("Some tests failed")
	}
}