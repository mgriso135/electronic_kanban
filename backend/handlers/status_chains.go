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

// GetStatusChainsHandler returns a handler for GET /api/status-chains
func GetStatusChainsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statusChains, err := getStatusChains(db)
		if err != nil {
			http.Error(w, "Failed to fetch status chains", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statusChains)
	}
}

// CreateStatusChainHandler returns a handler for POST /api/status-chains
func CreateStatusChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateStatusChainHandler: Starting") // Log entry
		var statusChainRequest struct {
			StatusChain     models.StatusChain         `json:"status_chain"`
			StatusesUpdates []models.StatusChainStatus `json:"statuses"` // Expecting statuses in request now
		}
		err := json.NewDecoder(r.Body).Decode(&statusChainRequest)
		if err != nil {
			log.Printf("CreateStatusChainHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("CreateStatusChainHandler: Request body decoded, data: %+v", statusChainRequest) // Log received data

		newStatusChain, err := createStatusChain(db, statusChainRequest.StatusChain)
		if err != nil {
			log.Printf("CreateStatusChainHandler: Failed to create status chain: %v", err)
			http.Error(w, "Failed to create status chain", http.StatusInternalServerError)
			return
		}
		log.Printf("CreateStatusChainHandler: Status chain created successfully, ID: %d", newStatusChain.StatusChainID) // Log success

		// Insert linked statuses into status_chains_statuses
		if len(statusChainRequest.StatusesUpdates) > 0 {
			err = insertStatusChainStatuses(db, newStatusChain.StatusChainID, statusChainRequest.StatusesUpdates)
			if err != nil {
				log.Printf("CreateStatusChainHandler: Failed to insert status chain statuses: %v", err)
				http.Error(w, "Failed to insert status chain statuses", http.StatusInternalServerError)
				return
			}
			log.Println("CreateStatusChainHandler: Status chain statuses inserted successfully")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newStatusChain)
		log.Println("CreateStatusChainHandler: Finished successfully") // Log exit
	}
}

// GetStatusChainHandler returns a handler for GET /api/status-chains/{id}
func GetStatusChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status chain ID", http.StatusBadRequest)
			return
		}

		statusChain, err := getStatusChainByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch status chain", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statusChain)
	}
}

// UpdateStatusChainHandler returns a handler for PUT/PATCH /api/status-chains/{id}
func UpdateStatusChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status chain ID", http.StatusBadRequest)
			return
		}

		var statusChainUpdates models.StatusChain
		err = json.NewDecoder(r.Body).Decode(&statusChainUpdates)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		statusChainUpdates.StatusChainID = id // Ensure ID from URL is used

		updatedStatusChain, err := updateStatusChain(db, statusChainUpdates)
		if err != nil {
			http.Error(w, "Failed to update status chain", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedStatusChain)
	}
}

// DeleteStatusChainHandler returns a handler for DELETE /api/status-chains/{id}
func DeleteStatusChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status chain ID", http.StatusBadRequest)
			return
		}

		err = deleteStatusChain(db, id)
		if err != nil {
			http.Error(w, "Failed to delete status chain", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
		json.NewEncoder(w).Encode(map[string]string{"message": "Status chain deleted"})
	}
}

// GetStatusChainStatusesHandler returns a handler to GET statuses for a specific chain
func GetStatusChainStatusesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		statusChainIDStr, ok := vars["statusChainId"]
		if !ok {
			http.Error(w, "Status Chain ID is required", http.StatusBadRequest)
			return
		}
		statusChainID, err := strconv.ParseInt(statusChainIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status chain ID", http.StatusBadRequest)
			return
		}

		statuses, err := getStatusChainStatuses(db, statusChainID)
		if err != nil {
			http.Error(w, "Failed to fetch statuses for status chain", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	}
}

// UpdateStatusChainStatusesHandler handles PUT requests to update statuses in a status chain (order, customer_supplier)
func UpdateStatusChainStatusesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateStatusChainStatusesHandler: Starting") // Log entry
		vars := mux.Vars(r)
		statusChainIDStr, ok := vars["statusChainId"]
		if !ok {
			http.Error(w, "Status Chain ID is required", http.StatusBadRequest)
			return
		}
		statusChainID, err := strconv.ParseInt(statusChainIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status chain ID", http.StatusBadRequest)
			return
		}

		var statusChainStatusesUpdates []models.StatusChainStatus // Assuming you'll create this model
		err = json.NewDecoder(r.Body).Decode(&statusChainStatusesUpdates)
		if err != nil {
			log.Printf("UpdateStatusChainStatusesHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("UpdateStatusChainStatusesHandler: Request body decoded, StatusChainID: %d, StatusUpdates: %+v", statusChainID, statusChainStatusesUpdates) // Log received data

		updatedStatuses, err := updateStatusChainStatuses(db, statusChainID, statusChainStatusesUpdates)
		if err != nil {
			log.Printf("UpdateStatusChainStatusesHandler: Failed to update statuses for status chain: %v", err)
			http.Error(w, "Failed to update statuses for status chain", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedStatuses)
		log.Println("UpdateStatusChainStatusesHandler: Finished successfully") // Log exit
	}
}

// Database interaction functions (private)

func getStatusChains(db *sql.DB) ([]models.StatusChain, error) {
	rows, err := db.Query("SELECT status_chain_id, name FROM status_chains")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statusChains []models.StatusChain
	for rows.Next() {
		var statusChain models.StatusChain
		if err := rows.Scan(&statusChain.StatusChainID, &statusChain.Name); err != nil {
			return nil, err
		}
		statusChains = append(statusChains, statusChain)
	}
	return statusChains, nil
}

func createStatusChain(db *sql.DB, statusChain models.StatusChain) (*models.StatusChain, error) {
	log.Printf("createStatusChain: Starting with data: %+v", statusChain) // Log entry
	sqlStatement := `
		INSERT INTO status_chains (name)
		VALUES ($1)
		RETURNING status_chain_id, name`
	log.Printf("createStatusChain: SQL Query: %s, Parameters: [%s]", sqlStatement, statusChain.Name) // Log SQL query
	var newStatusChain models.StatusChain
	err := db.QueryRow(sqlStatement, statusChain.Name).Scan(
		&newStatusChain.StatusChainID, &newStatusChain.Name,
	)
	if err != nil {
		log.Printf("createStatusChain: Error executing query: %v", err) // Log error
		return nil, fmt.Errorf("createStatusChain: error executing query: %w", err)
	}
	log.Printf("createStatusChain: Status chain created succesfully, ID: %d, Name: %s", newStatusChain.StatusChainID, newStatusChain.Name) // Log success

	return &newStatusChain, nil
}

func getStatusChainByID(db *sql.DB, id int64) (*models.StatusChain, error) {
	sqlStatement := `SELECT status_chain_id, name FROM status_chains WHERE status_chain_id = $1`
	var statusChain models.StatusChain
	err := db.QueryRow(sqlStatement, id).Scan(&statusChain.StatusChainID, &statusChain.Name)
	if err != nil {
		return nil, err
	}
	return &statusChain, nil
}

func updateStatusChain(db *sql.DB, statusChain models.StatusChain) (*models.StatusChain, error) {
	sqlStatement := `
		UPDATE status_chains
		SET name = $2
		WHERE status_chain_id = $1
		RETURNING status_chain_id, name`
	var updatedStatusChain models.StatusChain
	err := db.QueryRow(sqlStatement, statusChain.StatusChainID, statusChain.Name).Scan(
		&updatedStatusChain.StatusChainID, &updatedStatusChain.Name,
	)
	if err != nil {
		return nil, err
	}
	return &updatedStatusChain, nil
}

func deleteStatusChain(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM status_chains WHERE status_chain_id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}

// getStatusChainStatuses retrieves statuses for a given status chain, ordered by their 'order' field.
func getStatusChainStatuses(db *sql.DB, statusChainID int64) ([]map[string]interface{}, error) {
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
		return nil, fmt.Errorf("error querying status chain statuses: %w", err)
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
			return nil, fmt.Errorf("error scanning status chain status row: %w", err)
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
		return nil, fmt.Errorf("error iterating status chain statuses rows: %w", err)
	}

	return statuses, nil
}

// insertStatusChainStatuses inserts multiple status_chains_statuses records in a transaction
func insertStatusChainStatuses(db *sql.DB, statusChainID int64, statusesUpdates []models.StatusChainStatus) error {
	log.Println("insertStatusChainStatuses: Starting, StatusChainID:", statusChainID, ", StatusUpdates:", statusesUpdates) // Log entry
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("insertStatusChainStatuses: error starting transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO status_chains_statuses (status_chain_id, status_id, "order", customer_supplier)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("insertStatusChainStatuses: error preparing statement: %w", err)
	}
	defer stmt.Close()

	for _, statusUpdate := range statusesUpdates {
		log.Printf("updateStatusChainStatuses: Processing status update: StatusChainID=%d, StatusUpdate=%+v", statusChainID, statusUpdate) // Log each status update
		_, err = stmt.Exec(statusChainID, statusUpdate.StatusID, statusUpdate.Order, statusUpdate.CustomerSupplier)
		if err != nil {
			log.Printf("insertStatusChainStatuses: Error inserting status (StatusID: %d) for chain: %v", statusUpdate.StatusID, err) // Log individual status insert error
			return fmt.Errorf("insertStatusChainStatuses: error executing insert for status_id %d: %w", statusUpdate.StatusID, err)
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Printf("insertStatusChainStatuses: Error committing transaction: %v", err) // Log transaction commit error
		return fmt.Errorf("insertStatusChainStatuses: error committing transaction: %w", err)
	}
	log.Println("insertStatusChainStatuses: Status chain statuses inserted successfully") // Log success
	return nil
}

// updateStatusChainStatuses updates the statuses within a status chain (order, customer_supplier).
func updateStatusChainStatuses(db *sql.DB, statusChainID int64, statusesUpdates []models.StatusChainStatus) ([]map[string]interface{}, error) {
	log.Println("updateStatusChainStatuses: Starting, StatusChainID:", statusChainID, ", StatusUpdates:", statusesUpdates) // Log entry
	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("updateStatusChainStatuses: error starting transaction: %w", err)
	}
	defer tx.Rollback() // Rollback if function exits early

	// Prepare statement for updating status_chains_statuses
	stmt, err := tx.Prepare(`
		UPDATE status_chains_statuses
		SET "order" = $3, customer_supplier = $4
		WHERE status_chain_id = $1 AND status_id = $2
	`)
	if err != nil {
		return nil, fmt.Errorf("updateStatusChainStatuses: error preparing update statement: %w", err)
	}
	defer stmt.Close()

	for _, statusUpdate := range statusesUpdates {
		log.Printf("updateStatusChainStatuses: Processing status update: StatusChainID=%d, StatusUpdate=%+v", statusChainID, statusUpdate) // Log each status update
		_, err := stmt.Exec(statusChainID, statusUpdate.StatusID, statusUpdate.Order, statusUpdate.CustomerSupplier)
		if err != nil {
			log.Printf("updateStatusChainStatuses: Error updating status chain status (status_id: %d): %v", statusUpdate.StatusID, err) // Log error for each status update
			return nil, fmt.Errorf("updateStatusChainStatuses: error updating status chain status (status_id: %d): %w", statusUpdate.StatusID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		log.Printf("updateStatusChainStatuses: Error committing transaction: %v", err) // Log transaction commit error
		return nil, fmt.Errorf("updateStatusChainStatuses: error committing transaction: %w", err)
	}

	// After successful update, retrieve and return the updated statuses for the chain
	updatedStatuses, err := getStatusChainStatuses(db, statusChainID)
	if err != nil {
		log.Printf("updateStatusChainStatuses: Error retrieving updated status chain statuses: %v", err) // Log error on retrieval
		return nil, fmt.Errorf("error retrieving updated status chain statuses: %w", err)
	}
	log.Println("updateStatusChainStatuses: Status chain statuses updated successfully") // Log success

	return updatedStatuses, nil
}
