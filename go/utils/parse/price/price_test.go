package price

import (
	"testing"

	"github.com/jakubruminski/FYP/go/utils/logger"
)

func TestFloat(t *testing.T) {
	logger := &logger.Logger{}

	testCases := []struct {
		priceType string
		price     string
		expected  float64
	}{
		{"litres", "€43 per 70cl", 61.42857142857143},
		{"litres", "€50 per 100cl", 50.0},
		{"kilograms", "$5/kg", 5.0},
		{"kilograms", "€2/g", 2000.0},
		{"litres", "£3/ml", 3000.0},
		{"litres", "$0.01/cl", 1.0},
		{"each", "$2/item", 2.0},
		{"kilograms", "5/kg", 5.0},
		{"kilograms", "2/g", 2000.0},
		{"litres", "3/ml", 3000.0},
		{"litres", "0.01/cl", 1.0},
		{"each", "2/item", 2.0},
		{"litres", "€43 per 70cl", 61.42857142857143}, 

		// These should fail
		{"kilograms", "$5", 0.0},
		{"kilograms", "invalid", 0.0},
		{"unknown", "$5/kg", 0.0},

	}

	for i, tc := range testCases {

		logger.INFO("%d STARTING - PriceType: %s, Price: %s, Expected: %f", i, tc.priceType, tc.price, tc.expected)
		_, result, ok := Float(logger, tc.price)
		if tc.expected == 0.0 {
			if ok {
				t.Errorf("Expected ok to be false, but got true for test case PriceType: %s, Price: %s", tc.priceType, tc.price)
			}
		} else {
			if !ok {
				t.Errorf("Expected ok to be true, but got false for test case PriceType: %s, Price: %s", tc.priceType, tc.price)
			}
		}

		if result != tc.expected {
			t.Errorf("Expected %f, but got %f for test case PriceType: %s, Price: %s", tc.expected, result, tc.priceType, tc.price)
		}

		logger.INFO("%d FINISHED - PriceType: %s, Price: %s, Result: %f, Expected: %f", i, tc.priceType, tc.price, result, tc.expected)
	}
}