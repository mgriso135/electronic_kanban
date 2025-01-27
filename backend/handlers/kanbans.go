package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"electronic_kanban_backend/models"

	"github.com/gorilla/mux"
)

// GetKanbansHandler returns a handler for GET /api/kanbans
func GetKanbansHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kanbans, err := getKanbans(db)
		if err != nil {
			http.Error(w, "Failed to fetch kanbans", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kanbans)
	}
}

// GetKanbanHandler returns a handler for GET /api/kanbans/{id}
func GetKanbanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Kanban ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid kanban ID", http.StatusBadRequest)
			return
		}

		kanban, err := getKanbanByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch kanban", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kanban)
	}
}

// UpdateKanbanHandler returns a handler for PUT/PATCH /api/kanbans/{id} - Primarily for status updates
func UpdateKanbanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Kanban ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid kanban ID", http.StatusBadRequest)
			return
		}

		var kanbanUpdates map[string]interface{} // Allow partial updates
		err = json.NewDecoder(r.Body).Decode(&kanbanUpdates)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		updatedKanban, err := updateKanban(db, id, kanbanUpdates)
		if err != nil {
			http.Error(w, "Failed to update kanban", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanban)
	}
}

// Database interaction functions (private)

func getKanbans(db *sql.DB) ([]models.Kanban, error) {
	rows, err := db.Query(`
		SELECT
			id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id,
			status_chain_id, status_current, tipo_contenitore, quantity
		FROM kanbans
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kanbans []models.Kanban
	for rows.Next() {
		var k models.Kanban
		if err := rows.Scan(
			&k.ID, &k.DataAggiornamento, &k.LeadtimeDays, &k.IsActive, &k.KanbanChainID,
			&k.StatusChainID, &k.StatusCurrent, &k.TipoContenitore, &k.Quantity,
		); err != nil {
			return nil, err
		}
		kanbans = append(kanbans, k)
	}
	return kanbans, nil
}

func getKanbanByID(db *sql.DB, id int64) (*models.Kanban, error) {
	sqlStatement := `
		SELECT
			id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id,
			status_chain_id, status_current, tipo_contenitore, quantity
		FROM kanbans
		WHERE id = $1
	`
	var kanban models.Kanban
	err := db.QueryRow(sqlStatement, id).Scan(
		&kanban.ID, &kanban.DataAggiornamento, &kanban.LeadtimeDays, &kanban.IsActive, &kanban.KanbanChainID,
		&kanban.StatusChainID, &kanban.StatusCurrent, &kanban.TipoContenitore, &kanban.Quantity,
	)
	if err != nil {
		return nil, err
	}
	return &kanban, nil
}

func updateKanban(db *sql.DB, id int64, updates map[string]interface{}) (*models.Kanban, error) {
	// Start building the UPDATE query dynamically
	sqlStatement := `UPDATE kanbans SET data_aggiornamento = NOW()` // Always update data_aggiornamento
	var args []interface{}
	argIndex := 1

	if statusCurrent, ok := updates["status_current"].(float64); ok { // Assuming status_current is sent as integer in JSON, hence float64 after unmarshaling
		sqlStatement += fmt.Sprintf(", status_current = $%d", argIndex+1)
		args = append(args, int64(statusCurrent)) // Convert float64 to int64
		argIndex++
	}
	// Add other fields to update here if needed, following the same pattern

	sqlStatement += fmt.Sprintf(" WHERE id = $%d RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity", argIndex+1)
	args = append(args, id)

	var updatedKanban models.Kanban
	row := db.QueryRow(sqlStatement, args...) // Pass all arguments as slice
	err := row.Scan(
		&updatedKanban.ID, &updatedKanban.DataAggiornamento, &updatedKanban.LeadtimeDays, &updatedKanban.IsActive, &updatedKanban.KanbanChainID,
		&updatedKanban.StatusChainID, &updatedKanban.StatusCurrent, &updatedKanban.TipoContenitore, &updatedKanban.Quantity,
	)
	if err != nil {
		return nil, err
	}

	// Record history of status change if status_current was updated
	if _, statusUpdated := updates["status_current"]; statusUpdated {
		if err := recordKanbanHistory(db, &updatedKanban, updates); err != nil {
			// Log the history recording error, but don't fail the main update
			fmt.Printf("Error recording kanban history: %v\n", err) // Or use a proper logger
		}
	}

	return &updatedKanban, nil
}

func recordKanbanHistory(db *sql.DB, updatedKanban *models.Kanban, updates map[string]interface{}) error {
	var previousStatus int64
	if statusCurrentFloat, ok := updates["status_current"].(float64); ok {
		currentStatus := int64(statusCurrentFloat)
		// Retrieve previous status from the kanban record before update
		previousKanban, err := getKanbanByID(db, updatedKanban.ID) // Get Kanban before update to find previous status
		if err != nil {
			return fmt.Errorf("error fetching previous kanban status: %w", err)
		}
		previousStatus = previousKanban.StatusCurrent

		sqlStatement := `
			INSERT INTO kanban_histories (kanban_id, previous_status, next_status, data_aggiornamento)
			VALUES ($1, $2, $3, NOW())
		`
		_, err = db.Exec(sqlStatement, updatedKanban.ID, previousStatus, currentStatus)
		if err != nil {
			return fmt.Errorf("error inserting kanban history: %w", err)
		}
	}
	return nil
}
