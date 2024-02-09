package baskets


type User struct {
	UserID          int       `json:"user_id"`           // User's unique ID
	LastFetch       int       `json:"last_fetch"`        // Timestamp
	ProductID       int       `json:"product_id"`        // Reference to the product in the products table
	ProductExists   bool      `json:"product_exists"`    // This is a flag allowing to check if the product exists
}

