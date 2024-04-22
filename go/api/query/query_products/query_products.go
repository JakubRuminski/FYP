package query_products

import (
	"context"
	"database/sql"

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

func Get(logger *logger.Logger, tx *sql.Tx, products *[]*product.Product, productIDs *[]*int64) (ok bool) {

    query := `SELECT id, seller, name, currency, price, price_per_unit, discount_price, discount_price_per_unit, discount_price_in_words, unit_type, url, img_url FROM products WHERE id = ANY($1)`
    ok = postgres.ExecuteContextLookUpQuery(logger, tx, get, query, productIDs, products)
    if !ok {
        logger.ERROR("Failed to get products")
        return false
    }

	return true
}

func get(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

    productIDs, ok := args[0].(*[]*int64)
    if !ok {
        logger.ERROR("Failed to get product IDs")
        return false
    }

    rows, err := tx.QueryContext(ctx, query, pq.Array(*productIDs))
    if err != nil {
        logger.ERROR("Failed to get products: %s", err)
        return false
    }
    defer rows.Close()

    products, ok := args[1].(*[]*product.Product)

    i := 0
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
        
        logger.DEBUG("%d: Product: --------------------------", i)
        logger.DEBUG("%d: p.ID: %d", i, product.ID)
        logger.DEBUG("%d: p.Seller: %s", i, product.Seller)
        logger.DEBUG("%d: p.Name: %s", i, product.Name)
        logger.DEBUG("%d: p.Currency: %s", i, product.Currency)
        logger.DEBUG("%d: p.Price: %f", i, product.Price)
        logger.DEBUG("%d: p.PricePerUnit: %f", i, product.PricePerUnit)
        logger.DEBUG("%d: p.DiscountPrice: %f", i, product.DiscountPrice)
        logger.DEBUG("%d: p.DiscountPricePerUnit: %f", i, product.DiscountPricePerUnit)
        logger.DEBUG("%d: p.DiscountPriceInWords: %s", i, product.DiscountPriceInWords)
        logger.DEBUG("%d: p.UnitType: %s", i, product.UnitType)
        logger.DEBUG("%d: p.URL: %s", i, product.URL)
        logger.DEBUG("%d: p.ImgURL: %s", i, product.ImgURL)
        i++
    }

    return true
}


func Add(logger *logger.Logger, tx *sql.Tx, oldProducts, products *[]*product.Product) bool {

	// Check if there are products to insert
	if len(*products) == 0 {
		logger.INFO("No products to add")
		return true
	}

    query := product.ProductInsertQuery()

	ok := postgres.ExecuteContextChangeQuery(logger, tx, add, query, oldProducts, products)
	if !ok {
		logger.ERROR("Failed to add products")
		return false
	}

	return true
}


func add(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {

    products, ok := args[1].(*[]*product.Product)
    if !ok {
        logger.ERROR("Failed to get products")
        return false
    }

    // Accesses oldProducts
    _, ok = args[0].(*[]*product.Product)
    if !ok {
        logger.ERROR("Failed to get oldProducts")
        return false
    }

    nextAvailableID, ok := getNextAvailableID(logger, tx)
    if !ok {
        logger.ERROR("Failed to get the next available ID")
        return false
    }

    rowsAffected := int64(0)
    for _, product := range *products {

        // TODO: Possibly check oldProducts if match and just return it's ID.

        id, exists, ok := productUrlAlreadyExists(logger, tx, product.URL)
        if !ok {
            logger.ERROR("Failed to check if the product URL exists")
            return false
        }
        if exists {
            product.ID = id
            logger.DEBUG("Product URL already exists at ID: %d", product.ID)
            continue
        }

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


func productUrlAlreadyExists(logger *logger.Logger, tx *sql.Tx, url string) (id int64, exists, ok bool) {

    query := `SELECT id FROM products WHERE url = $1`
    ok = postgres.ExecuteContextLookUpQuery(logger, tx, existsQuery, query, url, &id, &exists)
    if !ok {
        logger.ERROR("Failed to check if the product URL exists")
        return -1, false, false
    }

    return id, exists, true
}

func existsQuery(logger *logger.Logger, tx *sql.Tx, ctx context.Context, query string, args ...interface{}) (ok bool) {
        
    url, ok := args[0].(string)
    if !ok {
        logger.ERROR("Failed to get the URL")
        return false
    }

    id, ok := args[1].(*int64)
    if !ok {
        logger.ERROR("Failed to get the ID")
        return false
    }

    exists, ok := args[2].(*bool)
    if !ok {
        logger.ERROR("Failed to get the exists")
        return false
    }

    err := tx.QueryRowContext(ctx, query, url).Scan(id)
    if err == sql.ErrNoRows {
        logger.DEBUG("Product URL doesn't exist")
        *exists = false
        return true
    }

    if err != nil {
        logger.ERROR("Failed to check if the product URL exists. Reason: %s", err)
        return false
    }

    logger.DEBUG("Product URL exists at ID: %d", *id)
    *exists = true
    return true
}


func getNextAvailableID(logger *logger.Logger, tx *sql.Tx) (id int64, ok bool) {

    query := `SELECT nextval('products_id_seq')`
    ok = postgres.ExecuteContextLookUpQuery(logger, tx, getID, query, &id)
    if !ok {
        logger.ERROR("Failed to get the next available ID")
        return -1, false
    }

    logger.DEBUG("Next available ID: %d", id)

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