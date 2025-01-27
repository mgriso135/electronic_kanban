package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"electronic_kanban_backend/models"

	"github.com/gorilla/mux"
)

// GetStatusesHandler returns a handler for GET /api/statuses
func GetStatusesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		statuses, err := getStatuses(db)
		if err != nil {
			http.Error(w, "Failed to fetch statuses", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(statuses)
	}
}

// CreateStatusHandler returns a handler for POST /api/statuses
func CreateStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var status models.Status
		err := json.NewDecoder(r.Body).Decode(&status)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newStatus, err := createStatus(db, status)
		if err != nil {
			http.Error(w, "Failed to create status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newStatus)
	}
}

// GetStatusHandler returns a handler for GET /api/statuses/{id}
func GetStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status ID", http.StatusBadRequest)
			return
		}

		status, err := getStatusByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch status", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(status)
	}
}

// UpdateStatusHandler returns a handler for PUT/PATCH /api/statuses/{id}
func UpdateStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status ID", http.StatusBadRequest)
			return
		}

		var statusUpdates models.Status
		err = json.NewDecoder(r.Body).Decode(&statusUpdates)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		statusUpdates.StatusID = id // Ensure ID from URL is used

		updatedStatus, err := updateStatus(db, statusUpdates)
		if err != nil {
			http.Error(w, "Failed to update status", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedStatus)
	}
}

// DeleteStatusHandler returns a handler for DELETE /api/statuses/{id}
func DeleteStatusHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Status ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid status ID", http.StatusBadRequest)
			return
		}

		err = deleteStatus(db, id)
		if err != nil {
			http.Error(w, "Failed to delete status", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
		json.NewEncoder(w).Encode(map[string]string{"message": "Status deleted"})
	}
}

// Database interaction functions (private)

func getStatuses(db *sql.DB) ([]models.Status, error) {
	rows, err := db.Query("SELECT status_id, name, color FROM statuses")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var statuses []models.Status
	for rows.Next() {
		var status models.Status
		if err := rows.Scan(&status.StatusID, &status.Name, &status.Color); err != nil {
			return nil, err
		}
		statuses = append(statuses, status)
	}
	return statuses, nil
}

func createStatus(db *sql.DB, status models.Status) (*models.Status, error) {
	sqlStatement := `
		INSERT INTO statuses (name, color)
		VALUES ($1, $2)
		RETURNING status_id, name, color`
	var newStatus models.Status
	err := db.QueryRow(sqlStatement, status.Name, status.Color).Scan(
		&newStatus.StatusID, &newStatus.Name, &newStatus.Color,
	)
	if err != nil {
		return nil, err
	}
	return &newStatus, nil
}

func getStatusByID(db *sql.DB, id int64) (*models.Status, error) {
	sqlStatement := `SELECT status_id, name, color FROM statuses WHERE status_id = $1`
	var status models.Status
	err := db.QueryRow(sqlStatement, id).Scan(&status.StatusID, &status.Name, &status.Color)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func updateStatus(db *sql.DB, status models.Status) (*models.Status, error) {
	sqlStatement := `
		UPDATE statuses
		SET name = $2, color = $3
		WHERE status_id = $1
		RETURNING status_id, name, color`
	var updatedStatus models.Status
	err := db.QueryRow(sqlStatement, status.StatusID, status.Name, status.Color).Scan(
		&updatedStatus.StatusID, &updatedStatus.Name, &updatedStatus.Color,
	)
	if err != nil {
		return nil, err
	}
	return &updatedStatus, nil
}

func deleteStatus(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM statuses WHERE status_id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}
