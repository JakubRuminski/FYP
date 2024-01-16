package product

type Product struct {
	Seller        string  `json:"seller"`
	ID            string  `json:"id"`              // product id. Database specific
	Name          string  `json:"name"`
	Price         float64 `json:"price"`
	SubPrice      float64 `json:"sub_price"`
	DiscountPrice float64 `json:"discount_price"`
	URL           string  `json:"url"`
	ImgURL        string  `json:"img_url"`
}