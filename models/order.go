package models

type Order struct {
	ID         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	ProductID  int     `json:"product_id"`
	Quantity   int     `json:"quantity"`
	Total      float64 `json:"total"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}
