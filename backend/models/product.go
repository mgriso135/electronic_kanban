package models

// Product model for products table
type Product struct {
	ProductID string `json:"product_id"`
	Name      string `json:"name"`
}
