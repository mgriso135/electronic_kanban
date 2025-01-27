package models

import "time"

// Kanban model for kanbans table
type Kanban struct {
	ID                int64     `json:"id"`
	DataAggiornamento time.Time `json:"data_aggiornamento"`
	LeadtimeDays      int64     `json:"leadtime_days"`
	IsActive          bool      `json:"is_active"`
	KanbanChainID     int64     `json:"kanban_chain_id"`
	StatusChainID     int64     `json:"status_chain_id"`
	StatusCurrent     int64     `json:"status_current"`
	TipoContenitore   string    `json:"tipo_contenitore"`
	Quantity          float64   `json:"quantity"`
}
