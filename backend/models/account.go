package models

// Account model for accounts table
type Account struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	VATNumber string `json:"vat_number"`
	Address   string `json:"address"`
}
