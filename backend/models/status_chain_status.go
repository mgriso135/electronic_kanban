package models

// StatusChainStatus model for status_chains_statuses table
type StatusChainStatus struct {
	StatusChainID    int64 `json:"status_chain_id"`
	StatusID         int64 `json:"status_id"`
	Order            int64 `json:"order"`
	CustomerSupplier int   `json:"customer_supplier"` // 1=Supplier, 2=Customer
}
