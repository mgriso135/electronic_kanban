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

// GetKanbanChainsHandler returns a handler for GET /api/kanban-chains
func GetKanbanChainsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		kanbanChains, err := getKanbanChains(db)
		if err != nil {
			http.Error(w, "Failed to fetch kanban chains", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kanbanChains)
	}
}

// CreateKanbanChainHandler returns a handler for POST /api/kanban-chains
func CreateKanbanChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var kanbanChainRequest struct { // Anonymous struct to handle request body
			KanbanChain        models.KanbanChain `json:"kanban_chain"`
			NoOfInitialKanbans int64              `json:"no_of_initial_kanbans"`
		}

		err := json.NewDecoder(r.Body).Decode(&kanbanChainRequest)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newKanbanChain, err := createKanbanChain(db, kanbanChainRequest.KanbanChain)
		if err != nil {
			http.Error(w, "Failed to create kanban chain", http.StatusInternalServerError)
			return
		}

		// Create initial kanbans based on no_of_active_kanbans
		if kanbanChainRequest.NoOfInitialKanbans > 0 {
			err = createInitialKanbans(db, newKanbanChain.ID, kanbanChainRequest.NoOfInitialKanbans, newKanbanChain.StatusChainID, newKanbanChain.LeadtimeDays, newKanbanChain.TipoContenitore, newKanbanChain.Quantity)
			if err != nil {
				// Consider logging the error and perhaps rolling back the kanban chain creation in a transaction for full rollback.
				http.Error(w, "Failed to create initial kanbans for the chain", http.StatusInternalServerError)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newKanbanChain)
	}
}

// GetKanbanChainHandler returns a handler for GET /api/kanban-chains/{id}
func GetKanbanChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Kanban Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid kanban chain ID", http.StatusBadRequest)
			return
		}

		kanbanChain, err := getKanbanChainByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch kanban chain", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(kanbanChain)
	}
}

// UpdateKanbanChainHandler returns a handler for PUT/PATCH /api/kanban-chains/{id}
func UpdateKanbanChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Kanban Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid kanban chain ID", http.StatusBadRequest)
			return
		}

		var kanbanChainUpdates models.KanbanChain
		err = json.NewDecoder(r.Body).Decode(&kanbanChainUpdates)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		kanbanChainUpdates.ID = id // Ensure ID from URL is used

		updatedKanbanChain, err := updateKanbanChain(db, kanbanChainUpdates)
		if err != nil {
			http.Error(w, "Failed to update kanban chain", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanbanChain)
	}
}

// DeleteKanbanChainHandler returns a handler for DELETE /api/kanban-chains/{id}
func DeleteKanbanChainHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			http.Error(w, "Kanban Chain ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid kanban chain ID", http.StatusBadRequest)
			return
		}

		err = deleteKanbanChain(db, id)
		if err != nil {
			http.Error(w, "Failed to delete kanban chain", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
		json.NewEncoder(w).Encode(map[string]string{"message": "Kanban chain deleted"})
	}
}

// Database interaction functions (private)

func getKanbanChains(db *sql.DB) ([]models.KanbanChain, error) {
	rows, err := db.Query(`
		SELECT
			id, cliente_id, prodotto_codice, fornitore_id, leadtime_days,
			quantity, tipo_contenitore, status_chain_id, no_of_active_kanbans
		FROM kanban_chains
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kanbanChains []models.KanbanChain
	for rows.Next() {
		var kc models.KanbanChain
		if err := rows.Scan(
			&kc.ID, &kc.ClienteID, &kc.ProdottoCodice, &kc.FornitoreID, &kc.LeadtimeDays,
			&kc.Quantity, &kc.TipoContenitore, &kc.StatusChainID, &kc.NoOfActiveKanbans,
		); err != nil {
			return nil, err
		}
		kanbanChains = append(kanbanChains, kc)
	}
	return kanbanChains, nil
}

func createKanbanChain(db *sql.DB, kc models.KanbanChain) (*models.KanbanChain, error) {
	sqlStatement := `
		INSERT INTO kanban_chains (
			cliente_id, prodotto_codice, fornitore_id, leadtime_days,
			quantity, tipo_contenitore, status_chain_id, no_of_active_kanbans
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING
			id, cliente_id, prodotto_codice, fornitore_id, leadtime_days,
			quantity, tipo_contenitore, status_chain_id, no_of_active_kanbans
	`
	var newKC models.KanbanChain
	err := db.QueryRow(sqlStatement,
		kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays,
		kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans,
	).Scan(
		&newKC.ID, &newKC.ClienteID, &newKC.ProdottoCodice, &newKC.FornitoreID, &newKC.LeadtimeDays,
		&newKC.Quantity, &newKC.TipoContenitore, &newKC.StatusChainID, &newKC.NoOfActiveKanbans,
	)
	if err != nil {
		return nil, err
	}
	return &newKC, nil
}

func getKanbanChainByID(db *sql.DB, id int64) (*models.KanbanChain, error) {
	sqlStatement := `
		SELECT
			id, cliente_id, prodotto_codice, fornitore_id, leadtime_days,
			quantity, tipo_contenitore, status_chain_id, no_of_active_kanbans
		FROM kanban_chains
		WHERE id = $1
	`
	var kc models.KanbanChain
	err := db.QueryRow(sqlStatement, id).Scan(
		&kc.ID, &kc.ClienteID, &kc.ProdottoCodice, &kc.FornitoreID, &kc.LeadtimeDays,
		&kc.Quantity, &kc.TipoContenitore, &kc.StatusChainID, &kc.NoOfActiveKanbans,
	)
	if err != nil {
		return nil, err
	}
	return &kc, nil
}

func updateKanbanChain(db *sql.DB, kc models.KanbanChain) (*models.KanbanChain, error) {
	sqlStatement := `
		UPDATE kanban_chains
		SET
			cliente_id = $2, prodotto_codice = $3, fornitore_id = $4, leadtime_days = $5,
			quantity = $6, tipo_contenitore = $7, status_chain_id = $8, no_of_active_kanbans = $9
		WHERE id = $1
		RETURNING
			id, cliente_id, prodotto_codice, fornitore_id, leadtime_days,
			quantity, tipo_contenitore, status_chain_id, no_of_active_kanbans
	`
	var updatedKC models.KanbanChain
	err := db.QueryRow(sqlStatement,
		kc.ID, kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays,
		kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans,
	).Scan(
		&updatedKC.ID, &updatedKC.ClienteID, &updatedKC.ProdottoCodice, &updatedKC.FornitoreID, &updatedKC.LeadtimeDays,
		&updatedKC.Quantity, &updatedKC.TipoContenitore, &updatedKC.StatusChainID, &updatedKC.NoOfActiveKanbans,
	)
	if err != nil {
		return nil, err
	}
	return &updatedKC, nil
}

func deleteKanbanChain(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM kanban_chains WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}

// createInitialKanbans creates kanban records when a new kanban chain is created
func createInitialKanbans(db *sql.DB, kanbanChainID int64, numberOfKanbans int64, statusChainID int64, leadtimeDays int64, tipoContenitore string, quantity float64) error {
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("error starting transaction for creating initial kanbans: %w", err)
	}
	defer tx.Rollback() // Rollback if any operation fails

	// Get the first status in the status chain to set as initial status_current
	firstStatusID, err := getFirstStatusIDInChain(db, statusChainID)
	if err != nil {
		return fmt.Errorf("error getting first status in chain: %w", err)
	}
	if firstStatusID == 0 {
		return fmt.Errorf("no statuses found in status chain %d", statusChainID)
	}

	sqlStatement := `
		INSERT INTO kanbans (
			kanban_chain_id, status_chain_id, status_current, leadtime_days, tipo_contenitore, quantity, data_aggiornamento
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	` // data_aggiornamento set to NOW() on creation

	stmt, err := tx.Prepare(sqlStatement)
	if err != nil {
		return fmt.Errorf("error preparing kanban insert statement: %w", err)
	}
	defer stmt.Close()

	for i := 0; i < int(numberOfKanbans); i++ {
		_, err = stmt.Exec(kanbanChainID, statusChainID, firstStatusID, leadtimeDays, tipoContenitore, quantity)
		if err != nil {
			return fmt.Errorf("error inserting kanban %d: %w", i+1, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transaction for initial kanban creation: %w", err)
	}

	return nil
}

// getFirstStatusIDInChain retrieves the status_id of the first status in a status chain (based on 'order').
func getFirstStatusIDInChain(db *sql.DB, statusChainID int64) (int64, error) {
	query := `
		SELECT status_id
		FROM status_chains_statuses
		WHERE status_chain_id = $1
		ORDER BY "order" ASC
		LIMIT 1;
	`
	var statusID int64
	err := db.QueryRow(query, statusChainID).Scan(&statusID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, nil // No statuses in chain, return 0
		}
		return 0, fmt.Errorf("error getting first status ID from chain: %w", err)
	}
	return statusID, nil
}
