package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
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
		log.Println("UpdateKanbanHandler: Starting")
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Println("UpdateKanbanHandler: Kanban ID is required")
			http.Error(w, "Kanban ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("UpdateKanbanHandler: Invalid kanban ID: %v", err)
			http.Error(w, "Invalid kanban ID", http.StatusBadRequest)
			return
		}

		var kanbanUpdates map[string]interface{} // Allow partial updates
		err = json.NewDecoder(r.Body).Decode(&kanbanUpdates)
		if err != nil {
			log.Printf("UpdateKanbanHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("UpdateKanbanHandler: Updating kanban with ID %d, data: %+v", id, kanbanUpdates)

		updatedKanban, err := updateKanbanStatus(db, id) // Call the new updateKanbanStatus function
		if err != nil {
			log.Printf("UpdateKanbanHandler: Failed to update kanban status: %v", err)
			http.Error(w, "Failed to update kanban status", http.StatusInternalServerError)
			return
		}

		log.Printf("UpdateKanbanHandler: Kanban updated successfully: %+v", updatedKanban)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanban)
		log.Println("UpdateKanbanHandler: Finished successfully")
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

// updateKanbanStatus updates the kanban status to the next status in the chain
func updateKanbanStatus(db *sql.DB, id int64) (*models.Kanban, error) {
	log.Printf("updateKanbanStatus: Starting, ID: %d", id)

	currentKanban, err := getKanbanByID(db, id)
	if err != nil {
		return nil, fmt.Errorf("updateKanbanStatus: error fetching kanban: %w", err)
	}

	statusChainStatuses, err := getStatusChainStatusesOrdered(db, currentKanban.StatusChainID)
	if err != nil {
		return nil, fmt.Errorf("updateKanbanStatus: error fetching status chain statuses: %w", err)
	}

	if len(statusChainStatuses) == 0 {
		return nil, fmt.Errorf("updateKanbanStatus: no statuses in status chain")
	}

	nextStatusID, err := getNextStatusInChain(currentKanban.StatusCurrent, statusChainStatuses)
	if err != nil {
		return nil, fmt.Errorf("updateKanbanStatus: error getting next status: %w", err)
	}

	sqlStatement := `
		UPDATE kanbans
		SET status_current = $2, data_aggiornamento = NOW()
		WHERE id = $1
		RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity
	`
	var updatedKanban models.Kanban
	err = db.QueryRow(sqlStatement, id, nextStatusID).Scan(
		&updatedKanban.ID, &updatedKanban.DataAggiornamento, &updatedKanban.LeadtimeDays, &updatedKanban.IsActive, &updatedKanban.KanbanChainID,
		&updatedKanban.StatusChainID, &updatedKanban.StatusCurrent, &updatedKanban.TipoContenitore, &updatedKanban.Quantity,
	)
	if err != nil {
		log.Printf("updateKanbanStatus: Error updating kanban status: %v", err)
		return nil, err
	}
	log.Printf("updateKanbanStatus: Kanban status updated succesfully to status_id: %d", nextStatusID)

	// Record history of status change
	updates := map[string]interface{}{"status_current": float64(nextStatusID)} // Pass nextStatusID as update
	if err := recordKanbanHistory(db, &updatedKanban, updates); err != nil {
		fmt.Printf("updateKanbanStatus: Error recording kanban history: %v\n", err)
	}

	return &updatedKanban, nil
}

// getStatusChainStatusesOrdered retrieves statuses for a given status chain, ordered by their 'order' field.
// It returns []map[string]interface{} to include 'order' and 'customer_supplier' if needed in the future.
func getStatusChainStatusesOrdered(db *sql.DB, statusChainID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT
			scs.status_id,
			s.name AS status_name,
			s.color AS status_color,
			scs."order",
			scs.customer_supplier
		FROM
			status_chains_statuses scs
		JOIN
			statuses s ON scs.status_id = s.status_id
		WHERE
			scs.status_chain_id = $1
		ORDER BY
			scs."order" ASC;
	`

	rows, err := db.Query(query, statusChainID)
	if err != nil {
		return nil, fmt.Errorf("getStatusChainStatusesOrdered: error querying status chain statuses: %w", err)
	}
	defer rows.Close()

	var statuses []map[string]interface{}
	for rows.Next() {
		var statusID int64
		var statusName string
		var statusColor string
		var order int64
		var customerSupplier int

		if err := rows.Scan(&statusID, &statusName, &statusColor, &order, &customerSupplier); err != nil {
			return nil, fmt.Errorf("getStatusChainStatusesOrdered: error scanning status chain status row: %w", err)
		}

		statuses = append(statuses, map[string]interface{}{
			"status_id":         statusID,
			"status_name":       statusName,
			"status_color":      statusColor,
			"order":             order,
			"customer_supplier": customerSupplier,
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getStatusChainStatusesOrdered: error iterating status chain statuses rows: %w", err)
	}

	return statuses, nil
}

// getNextStatusInChain determines the next status in the status chain based on the current status and chain order.
func getNextStatusInChain(currentStatusID int64, statusChainStatuses []map[string]interface{}) (int64, error) {
	var nextStatusID int64

	if len(statusChainStatuses) == 0 {
		return 0, fmt.Errorf("getNextStatusInChain: status chain is empty")
	}

	foundCurrent := false
	isLastStatus := false

	for i, statusMap := range statusChainStatuses {
		statusID := statusMap["status_id"].(int64)
		if statusID == currentStatusID {
			foundCurrent = true
			if i == len(statusChainStatuses)-1 {
				isLastStatus = true // Current status is the last one
			} else {
				nextStatusID = statusChainStatuses[i+1]["status_id"].(int64) // Get the next status
			}
			break // Exit loop once current status is found
		}
	}

	if !foundCurrent {
		return 0, fmt.Errorf("getNextStatusInChain: current status not found in status chain")
	}

	if isLastStatus {
		nextStatusID = statusChainStatuses[0]["status_id"].(int64) // Cycle back to the first status
	}

	return nextStatusID, nil
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
