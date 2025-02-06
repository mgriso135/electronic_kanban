package models

// KanbanChain model for kanban_chains table
type KanbanChain struct {
	ID                int64   `json:"id"`
	ClienteID         int64   `json:"cliente_id"`
	ProdottoCodice    string  `json:"prodotto_codice"`
	FornitoreID       int64   `json:"fornitore_id"`
	LeadtimeDays      int64   `json:"leadtime_days"`
	Quantity          float64 `json:"quantity"`
	TipoContenitore   string  `json:"tipo_contenitore"`
	StatusChainID     int64   `json:"status_chain_id"`
	NoOfActiveKanbans int64   `json:"no_of_active_kanbans"`
}
