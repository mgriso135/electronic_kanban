package models

// Status model for statuses table
type Status struct {
	StatusID int64  `json:"status_id"`
	Name     string `json:"name"`
	Color    string `json:"color"`
}
