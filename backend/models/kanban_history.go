package models

import "time"

// KanbanHistory model for kanban_histories table
type KanbanHistory struct {
	ID                int64     `json:"id"`
	KanbanID          int64     `json:"kanban_id"`
	PreviousStatus    int64     `json:"previous_status"`
	NextStatus        int64     `json:"next_status"`
	DataAggiornamento time.Time `json:"data_aggiornamento"`
}
