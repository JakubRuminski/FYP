package tesco

import (
	"testing"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/jakubruminski/FYP/go/utils/env"
	"github.com/jakubruminski/FYP/go/utils/logger"
)

func TestTesco(t *testing.T) {
	logger := &logger.Logger{}
	environment, ok := env.Get(logger, "ENVIRONMENT")
	if !ok { return }
	logger.Environment = environment

	parsedProducts := [][]string{
        {"tesco", "ID", "Tesco Fresh Milk 2 Litre", "€2.09", "€1.04/litre", "https://www.tesco.ie/groceries/en-IE/products/250005606", "https://digitalcontent.api.tesco.com/v2/media/ghs/1cf822bb-7dae-431e-b900-7be39d965d1d/6a7581fd-06d1-4d39-9aba-397eb74d8e0f_966300345.jpeg?h=225&w=225"}, 
		{"tesco", "ID", "Tesco Full Fat Milk 3Ltr", "€2.95", "€0.98/litre", "https://www.tesco.ie/groceries/en-IE/products/260776455", "https://digitalcontent.api.tesco.com/v2/media/ghs/28f7856e-97b1-4edb-a1b9-d6eb59a53604/snapshotimagehandler_417791735.jpeg?h=225&w=225"},

		
	}

	for _, parsedProduct := range parsedProducts {
		_, ok := product.NewProduct(logger, parsedProduct[0], parsedProduct[1], parsedProduct[2], parsedProduct[3], parsedProduct[4], parsedProduct[5], parsedProduct[6], parsedProduct[7], parsedProduct[8])
		if !ok {
			t.Errorf("Failed to create product using name %s, price %s, subPrice %s, specialPrice %s, link %s, imageURL %s", parsedProduct[2], parsedProduct[3], parsedProduct[4], parsedProduct[5], parsedProduct[6], parsedProduct[7])
		}
	}
}