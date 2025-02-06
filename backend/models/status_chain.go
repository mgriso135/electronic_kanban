package models

// StatusChain model for status_chains table
type StatusChain struct {
	StatusChainID int64  `json:"status_chain_id"`
	Name          string `json:"name"`
}
