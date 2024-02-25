package query_products

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/jakubruminski/FYP/go/api/product"
	"github.com/lib/pq"

	"github.com/jakubruminski/FYP/go/utils/logger"
	"github.com/jakubruminski/FYP/go/utils/postgres"
)

var tableName = "products"

func INIT(logger *logger.Logger) (ok bool) {

	query := product.ProductCreateQuery()

	ok = postgres.ExecuteCreateTableQuery(logger, tableName, query)
	if !ok {
		logger.ERROR("Couldn't create the products table")
		return false
	}

	return true
}

func Get(logger *logger.Logger, tx *sql.Tx, products *[]*product.Product, productIDs *[]*string) (ok bool) {

    query := `SELECT id, seller, name, currency, price, price_per_unit, discount_price, discount_price_per_unit, discount_price_in_words, unit_type, url, img_url FROM products WHERE id = ANY($1)`
    ok = postgres.ExecuteContextLookUpQuery(logger, tx, get, query, productIDs, products)
    if !ok {
        logger.ERROR("Failed to get products")
        return false
    }

	return true
}

func get(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

    productIDs, ok := args[0].(*[]*string)
    if !ok {
        logger.ERROR("Failed to get product IDs")
        return false
    }

    var productIDsSlice []int64
    for _, productID := range *productIDs {
        id, err := strconv.ParseInt(*productID, 10, 64)
        if err != nil {
            logger.ERROR("Failed to parse product ID: %s", err)
            return false
        }
        productIDsSlice = append(productIDsSlice, id)
    }

    rows, err := tx.QueryContext(ctx, query, pq.Array(productIDsSlice))
    if err != nil {
        logger.ERROR("Failed to get products: %s", err)
        return false
    }
    defer rows.Close()

    products, ok := args[1].(*[]*product.Product)

    for rows.Next() {
        product := &product.Product{}
        err := rows.Scan(
            &product.ID,
            &product.Seller,
            &product.Name,
            &product.Currency,
            &product.Price,
            &product.PricePerUnit,
            &product.DiscountPrice,
            &product.DiscountPricePerUnit,
            &product.DiscountPriceInWords,
            &product.UnitType,
            &product.URL,
            &product.ImgURL,
        )
        if err != nil {
            logger.ERROR("Failed to scan product: %s", err)
            return false
        }
        *products = append(*products, product)
    }

    return true
}


func Add(logger *logger.Logger, tx *sql.Tx, products *[]*product.Product) bool {

	// Check if there are products to insert
	if len(*products) == 0 {
		logger.INFO("No products to add")
		return true
	}

    query := product.ProductInsertQuery()

	ok := postgres.ExecuteContextChangeQuery(logger, tx, add, query, products)
	if !ok {
		logger.ERROR("Failed to add products")
		return false
	}

	return true
}


func add(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

    products, ok := args[0].(*[]*product.Product)
    if !ok {
        logger.ERROR("Failed to get products")
        return false
    }

    nextAvailableID, ok := getNextAvailableID(logger, tx)
    if !ok {
        logger.ERROR("Failed to get the next available ID")
        return false
    }

    rowsAffected := int64(0)
    for _, product := range *products {
        result, err := tx.ExecContext(
            ctx,
            query,
            product.Seller,
            product.Name,
            product.Currency,
            product.Price,
            product.PricePerUnit,
            product.DiscountPrice,
            product.DiscountPricePerUnit,
            product.DiscountPriceInWords,
            product.UnitType,
            product.URL,
            product.ImgURL,
        )
        if err != nil {
            logger.ERROR("Failed to execute the query. Reason: %s", err)
            return false
        }
        rowsAffected, err = result.RowsAffected()
        if err != nil {
            logger.ERROR("Failed to get the number of rows affected. Reason: %s", err)
            return false
        }

        if rowsAffected != 1 {
            logger.ERROR("Expected 1 row to be affected, got %d", rowsAffected)
            return false
        }

        product.ID = nextAvailableID
        logger.DEBUG("Added product to the database")
        logger.DEBUG("ID: %d", product.ID)
        logger.DEBUG("Seller: %s", product.Seller)
        logger.DEBUG("Name: %s", product.Name)
        logger.DEBUG("Currency: %s", product.Currency)
        logger.DEBUG("Price: %f", product.Price)
        logger.DEBUG("PricePerUnit: %f", product.PricePerUnit)
        logger.DEBUG("DiscountPrice: %f", product.DiscountPrice)
        logger.DEBUG("DiscountPricePerUnit: %f", product.DiscountPricePerUnit)
        logger.DEBUG("DiscountPriceInWords: %s", product.DiscountPriceInWords)
        logger.DEBUG("UnitType: %s", product.UnitType)
        logger.DEBUG("URL: %s", product.URL)
        logger.DEBUG("ImgURL: %s", product.ImgURL)
        
        nextAvailableID++
    }

    return true
}


func getNextAvailableID(logger *logger.Logger, tx *sql.Tx) (id int64, ok bool) {

    query := `SELECT nextval('products_id_seq')`
    ok = postgres.ExecuteContextLookUpQuery(logger, tx, getID, query, &id)
    if !ok {
        logger.ERROR("Failed to get the next available ID")
        return -1, false
    }

    return id, true
}

func getID(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {
    
        var id int64
        err := tx.QueryRowContext(ctx, query).Scan(&id)
        if err != nil {
            logger.ERROR("Failed to get the next available ID. Reason: %s", err)
            return false
        }
    
        idPtr, ok := args[0].(*int64)
        if !ok {
            logger.ERROR("Failed to get the next available ID")
            return false
        }
        *idPtr = id
    
        return true
    }