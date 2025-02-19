package handlers

import (
	"database/sql"
	"fmt"

	"encoding/json"
	"log"
	"net/http"
	"strconv"

	// Import strings package
	"electronic_kanban_backend/models"

	"github.com/gorilla/mux"
	"github.com/lib/pq"
)

// GetKanbansHandler returns a handler for GET /api/kanbans, now with product filtering
func GetKanbansHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		productIDs := r.URL.Query()["product_id"] // Get product_id query parameters

		kanbans, err := getKanbans(db, productIDs) // Pass productIDs to getKanbans
		if err != nil {
			http.Error(w, "Failed to fetch kanbans", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kanbans)
	}
}

// CreateKanbanHandler returns a handler for POST /api/kanbans
func CreateKanbanHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var kanban models.Kanban
		err := json.NewDecoder(r.Body).Decode(&kanban)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newKanban, err := createKanban(db, kanban)
		if err != nil {
			http.Error(w, "Failed to create kanban", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newKanban)
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

		kanban, err := getKanbanByID(db, id) // Use getKanbanByID to fetch single kanban
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

// UpdateKanbanHandler returns a handler for PUT/PATCH /api/kanbans/{id} - For updating editable fields
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

		//		updatedKanban, err := updateKanbanPartial(db, id, kanbanUpdates) // Call new updateKanbanPartial function
		updatedKanban, err := updateKanbanStatus(db, id) // CALL CORRECT FUNCTION HERE - updateKanbanStatus

		if err != nil {
			log.Printf("UpdateKanbanHandler: Failed to update kanban: %v", err)
			http.Error(w, "Failed to update kanban", http.StatusInternalServerError)
			return
		}

		log.Printf("UpdateKanbanHandler: Kanban updated successfully: %+v", updatedKanban)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanban)
		log.Println("UpdateKanbanHandler: Finished successfully")
	}
}

// updateKanbanPartial updates specific fields of a kanban (leadtime_days, tipo_contenitore, quantity)
func updateKanbanPartial(db *sql.DB, id int64, updates map[string]interface{}) (*models.Kanban, error) {
	log.Printf("updateKanbanPartial: Starting, ID: %d, updates: %+v", id, updates)
	// Start building the UPDATE query dynamically
	sqlStatement := `UPDATE kanbans SET data_aggiornamento = NOW()` // Always update data_aggiornamento
	var args []interface{}
	argIndex := 1

	if leadtimeDays, ok := updates["leadtime_days"].(float64); ok {
		sqlStatement += fmt.Sprintf(", leadtime_days = $%d", argIndex)
		args = append(args, int64(leadtimeDays))
		argIndex++
	}
	if tipoContenitore, ok := updates["tipo_contenitore"].(string); ok {
		sqlStatement += fmt.Sprintf(", tipo_contenitore = $%d", argIndex)
		args = append(args, tipoContenitore)
		argIndex++
	}
	if quantity, ok := updates["quantity"].(float64); ok {
		sqlStatement += fmt.Sprintf(", quantity = $%d", argIndex)
		args = append(args, quantity)
		argIndex++
	}
	// Add other fields to update here if needed, following the same pattern

	sqlStatement += fmt.Sprintf(" WHERE id = $%d RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity", argIndex)
	args = append(args, id)

	log.Printf("updateKanbanPartial: SQL Query: %s, Parameters: %+v", sqlStatement, args)

	var updatedKanban models.Kanban
	row := db.QueryRow(sqlStatement, args...)
	err := row.Scan(
		&updatedKanban.ID, &updatedKanban.DataAggiornamento, &updatedKanban.LeadtimeDays, &updatedKanban.IsActive, &updatedKanban.KanbanChainID,
		&updatedKanban.StatusChainID, &updatedKanban.StatusCurrent, &updatedKanban.TipoContenitore, &updatedKanban.Quantity,
	)
	if err != nil {
		log.Printf("updateKanbanPartial: Error updating kanban: %v", err)
		return nil, err
	}
	log.Printf("updateKanbanPartial: Kanban updated succesfully with data: %+v", updatedKanban)
	return &updatedKanban, nil
}

// DeleteKanbanHandler returns a handler for DELETE /api/kanbans/{id}
func DeleteKanbanHandler(db *sql.DB) http.HandlerFunc {
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

		err = deleteKanban(db, id) // Call deleteKanban database function
		if err != nil {
			http.Error(w, "Failed to delete kanban", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
		json.NewEncoder(w).Encode(map[string]string{"message": "Kanban deleted"})
	}
}

// Database interaction functions (private)

// getKanbans retrieves kanbans, optionally filtered by product IDs
func getKanbans(db *sql.DB, productIDs []string) ([]map[string]interface{}, error) { // Return []map[string]interface{}
	query := `
		SELECT
			k.id AS kanban_id,
			k.data_aggiornamento,
			k.leadtime_days,
			k.is_active,
			k.kanban_chain_id,
			k.status_chain_id,
			k.status_current,
			k.tipo_contenitore,
			k.quantity,
			p.product_id AS product_id,  -- Include product_id
			p.name AS product_name       -- Include product_name
		FROM
			kanbans k
		JOIN
			kanban_chains kc ON k.kanban_chain_id = kc.id
		JOIN
			products p ON kc.prodotto_codice = p.product_id
		WHERE k.is_active=true
	`
	var args []interface{}
	if len(productIDs) > 0 {
		query += " AND kc.prodotto_codice = ANY($1)" // Filter directly on kc.prodotto_codice
		// Use pq.Array to correctly format the productIDs array for PostgreSQL
		args = append(args, pq.Array(productIDs)) // Use pq.Array here
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kanbans []map[string]interface{} // Return slice of maps
	for rows.Next() {
		var kanbanData map[string]interface{} = make(map[string]interface{}) // Use map[string]interface{}
		var k models.Kanban
		var productID string   // Temporary variable for product_id
		var productName string // Temporary variable for product_name

		err := rows.Scan(
			&k.ID, &k.DataAggiornamento, &k.LeadtimeDays, &k.IsActive, &k.KanbanChainID,
			&k.StatusChainID, &k.StatusCurrent, &k.TipoContenitore, &k.Quantity,
			&productID,   // Scan into temporary productID variable
			&productName, // Scan into temporary productName variable
		)
		if err != nil {
			return nil, err
		}
		// Assign temporary variables to the map
		kanbanData["product_id"] = productID
		kanbanData["product_name"] = productName

		// Copy Kanban fields to the map as well, if needed for frontend
		kanbanData["id"] = k.ID
		kanbanData["data_aggiornamento"] = k.DataAggiornamento
		kanbanData["leadtime_days"] = k.LeadtimeDays
		kanbanData["is_active"] = k.IsActive
		kanbanData["kanban_chain_id"] = k.KanbanChainID
		kanbanData["status_chain_id"] = k.StatusChainID
		kanbanData["status_current"] = k.StatusCurrent
		kanbanData["tipo_contenitore"] = k.TipoContenitore
		kanbanData["quantity"] = k.Quantity

		kanbans = append(kanbans, kanbanData) // Append the map to kanbans slice
	}
	return kanbans, nil
}

func createKanban(db *sql.DB, kanban models.Kanban) (*models.Kanban, error) {
	sqlStatement := `
		INSERT INTO kanbans (data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity)
		VALUES (NOW(), $1, $2, $3, $4, $5, $6, $7)
		RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity`
	var newKanban models.Kanban
	err := db.QueryRow(sqlStatement,
		kanban.LeadtimeDays, kanban.IsActive, kanban.KanbanChainID, kanban.StatusChainID, kanban.StatusCurrent, kanban.TipoContenitore, kanban.Quantity,
	).Scan(
		&newKanban.ID, &newKanban.DataAggiornamento, &newKanban.LeadtimeDays, &newKanban.IsActive, &newKanban.KanbanChainID,
		&newKanban.StatusChainID, &newKanban.StatusCurrent, &newKanban.TipoContenitore, &newKanban.Quantity,
	)
	if err != nil {
		return nil, err
	}
	return &newKanban, nil
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
	) // Line 289 is likely this rows.Scan line
	if err != nil {
		return nil, err
	}
	return &kanban, nil
}

// deleteKanban performs a soft delete of a kanban by setting is_active to false
func deleteKanban(db *sql.DB, id int64) error {
	sqlStatement := `
		UPDATE kanbans
		SET is_active = false, data_aggiornamento = NOW()
		WHERE id = $1
		RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity` // Returning updated kanban
	var updatedKanban models.Kanban // To scan the updated kanban
	err := db.QueryRow(sqlStatement, id).Scan(
		&updatedKanban.ID,
		&updatedKanban.DataAggiornamento,
		&updatedKanban.LeadtimeDays,
		&updatedKanban.IsActive,
		&updatedKanban.KanbanChainID,
		&updatedKanban.StatusChainID,
		&updatedKanban.StatusCurrent,
		&updatedKanban.TipoContenitore,
		&updatedKanban.Quantity, // **CORRECTED - Now exactly 9 arguments, matching RETURNING clause**
	)
	if err != nil {
		log.Printf("deleteKanban: Error executing soft DELETE (UPDATE is_active=false) query for kanban ID %d: %v", id, err)
		return fmt.Errorf("deleteKanban: error soft deleting kanban ID %d: %w", id, err)
	}
	log.Printf("deleteKanban: Successfully soft deleted kanban with ID %d", id)
	return nil
}

func updateKanban(db *sql.DB, id int64, updates map[string]interface{}) (*models.Kanban, error) {
	log.Printf("updateKanban: Starting, ID: %d, updates: %+v", id, updates)
	// Start building the UPDATE query dynamically
	sqlStatement := `UPDATE kanbans SET data_aggiornamento = NOW()` // Always update data_aggiornamento
	var args []interface{}
	argIndex := 1

	if statusCurrent, ok := updates["status_current"].(float64); ok { // Assuming status_current is sent as integer in JSON, hence float64 after unmarshaling
		sqlStatement += fmt.Sprintf(", status_current = $%d", argIndex)
		args = append(args, int64(statusCurrent)) // Convert float64 to int64
		argIndex++
	}
	// Add other fields to update here if needed, following the same pattern

	sqlStatement += fmt.Sprintf(" WHERE id = $%d RETURNING id, data_aggiornamento, leadtime_days, is_active, kanban_chain_id, status_chain_id, status_current, tipo_contenitore, quantity", argIndex)
	args = append(args, id)
	log.Printf("updateKanban: SQL Query: %s, Parameters: %+v", sqlStatement, args)

	var updatedKanban models.Kanban
	row := db.QueryRow(sqlStatement, args...) // Pass all arguments as slice
	err := row.Scan(
		&updatedKanban.ID, &updatedKanban.DataAggiornamento, &updatedKanban.LeadtimeDays, &updatedKanban.IsActive, &updatedKanban.KanbanChainID,
		&updatedKanban.StatusChainID, &updatedKanban.StatusCurrent, &updatedKanban.TipoContenitore, &updatedKanban.Quantity,
	)
	if err != nil {
		log.Printf("updateKanban: Error updating kanban: %v", err)
		return nil, err
	}
	log.Printf("updateKanban: Kanban updated succesfully with data: %+v", updatedKanban)

	// Record history of status change if status_current was updated
	if _, statusUpdated := updates["status_current"]; statusUpdated {
		log.Println("updateKanban: Recording status history")
		if err := recordKanbanHistory(db, &updatedKanban, updates); err != nil {
			// Log the history recording error, but don't fail the main update
			fmt.Printf("updateKanban: Error recording kanban history: %v\n", err) // Or use a proper logger
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

	nextStatusID, err := getNextStatusIDInChain(db, currentKanban.StatusCurrent, statusChainStatuses) // CALL getNextStatusIDInChain HERE
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
		log.Printf("updateKanbanStatus: Error updating kanban: %v", err)
		return nil, err
	}
	log.Printf("updateKanbanStatus: Kanban status updated succesfully to status_id: %d", nextStatusID)

	// Record history of status change if status_current was updated
	updates := map[string]interface{}{"status_current": float64(nextStatusID)} // Pass nextStatusID as update
	if err := recordKanbanHistory(db, &updatedKanban, updates); err != nil {
		fmt.Printf("updateKanbanStatus: Error recording kanban history: %v\n", err) // Or use a proper logger
	}

	return &updatedKanban, nil
}

// KanbanEditFormHandler handles PUT/PATCH requests to /api/kanbans/{id} - For Kanban Edit Form updates
func KanbanEditFormHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("KanbanEditFormHandler: Starting")
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Println("KanbanEditFormHandler: Kanban ID is required")
			http.Error(w, "Kanban ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("KanbanEditFormHandler: Invalid kanban ID: %v", err)
			http.Error(w, "Invalid kanban ID", http.StatusBadRequest)
			return
		}

		var kanbanUpdates map[string]interface{} // Allow partial updates
		err = json.NewDecoder(r.Body).Decode(&kanbanUpdates)
		if err != nil {
			log.Printf("KanbanEditFormHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("KanbanEditFormHandler: Updating kanban with ID %d, data from edit form: %+v", id, kanbanUpdates)

		updatedKanban, err := updateKanbanPartial(db, id, kanbanUpdates) // CALL updateKanbanPartial HERE - For Edit Form Updates
		if err != nil {
			log.Printf("KanbanEditFormHandler: Failed to update kanban: %v", err)
			http.Error(w, "Failed to update kanban", http.StatusInternalServerError)
			return
		}

		log.Printf("KanbanEditFormHandler: Kanban updated successfully via edit form: %+v", updatedKanban)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanban)
		log.Println("KanbanEditFormHandler: Finished successfully")
	}
}

// getNextStatusInChain determines the next status in the status chain based on the current status and chain order.
func getNextStatusIDInChain(db *sql.DB, currentStatusID int64, statusChainStatuses []map[string]interface{}) (int64, error) {
	log.Println("getNextStatusIDInChain: Starting - Current Status ID:", currentStatusID)          // LOG: Function start and currentStatusID
	log.Printf("getNextStatusIDInChain: Status Chain Statuses received: %+v", statusChainStatuses) // LOG: statusChainStatuses array

	var nextStatusID int64
	lastStatusID := int64(0) // To track the last status in the chain

	if len(statusChainStatuses) == 0 {
		log.Println("getNextStatusIDInChain: Status chain is EMPTY") // LOG: Empty chain
		return 0, fmt.Errorf("getNextStatusInChain: status chain is empty")
	}

	foundCurrent := false
	isLastStatus := false

	for i, statusMap := range statusChainStatuses {
		statusID := statusMap["status_id"].(int64)
		log.Printf("getNextStatusIDInChain: Checking status in chain - Index: %d, Status ID: %d, Order: %d", i, statusID, statusMap["order"]) // LOG: Status details in loop
		if i == len(statusChainStatuses)-1 {
			lastStatusID = statusID                                                         // Capture the last status ID
			log.Println("getNextStatusChain: Last Status ID in chain found:", lastStatusID) // Log: Last status ID
		}
		if statusID == currentStatusID {
			foundCurrent = true
			log.Println("getNextStatusInChain: Current Status ID FOUND in chain at index:", i) // Log: Current status found
			if i == len(statusChainStatuses)-1 {
				isLastStatus = true                                                             // Current status is the last one
				log.Println("getNextStatusInChain: Current Status is the LAST status in chain") // Log: Last status in chain
			} else {
				nextStatusID = statusChainStatuses[i+1]["status_id"].(int64)                  // Get the next status
				log.Println("getNextStatusInChain: Next Status ID determined:", nextStatusID) // Log: Next status ID
			}
			break // Exit loop once current status is found
		}
	}

	if !foundCurrent {
		log.Printf("getNextStatusIDInChain: Current Status ID %d NOT FOUND in status chain", currentStatusID) // Log: Current status NOT found
		return 0, fmt.Errorf("getNextStatusInChain: current status not found in status chain")
	}

	if isLastStatus {
		nextStatusID = statusChainStatuses[0]["status_id"].(int64)                          // Cycle back to the first status
		log.Println("getNextStatusInChain: Cycling back to FIRST Status ID:", nextStatusID) // Log: Cycle to first status
	}

	log.Println("getNextStatusInChain: Returning Next Status ID:", nextStatusID) // Log: Function exit and nextStatusID
	return nextStatusID, nil
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
