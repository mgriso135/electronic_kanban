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
		log.Println("CreateKanbanChainHandler: Starting") // Log entry

		var kanbanChainRequest struct { // Anonymous struct to handle request body
			KanbanChain        models.KanbanChain `json:"kanban_chain"`
			NoOfInitialKanbans int64              `json:"no_of_initial_kanbans"`
		}

		err := json.NewDecoder(r.Body).Decode(&kanbanChainRequest)
		if err != nil {
			log.Printf("CreateKanbanChainHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("CreateKanbanChainHandler: Request body decoded, Data: %+v", kanbanChainRequest)

		newKanbanChain, err := createKanbanChain(db, kanbanChainRequest.KanbanChain)
		if err != nil {
			log.Printf("CreateKanbanChainHandler: Failed to create kanban chain: %v", err)
			http.Error(w, "Failed to create kanban chain", http.StatusInternalServerError)
			return
		}

		log.Printf("CreateKanbanChainHandler: Kanban chain created successfully, ID: %d", newKanbanChain.ID)

		// Create initial kanbans based on no_of_active_kanbans
		if kanbanChainRequest.NoOfInitialKanbans > 0 {
			err = createInitialKanbans(db, newKanbanChain.ID, kanbanChainRequest.NoOfInitialKanbans, newKanbanChain.StatusChainID, newKanbanChain.LeadtimeDays, newKanbanChain.TipoContenitore, newKanbanChain.Quantity)
			if err != nil {
				// Consider logging the error and perhaps rolling back the kanban chain creation in a transaction for full rollback.
				http.Error(w, "Failed to create initial kanbans for the chain", http.StatusInternalServerError)
				log.Printf("CreateKanbanChainHandler: Failed to create initial kanbans for the chain: %v", err)
				return
			}
			log.Println("CreateKanbanChainHandler: initial kanbans created successfully")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newKanbanChain)
		log.Println("CreateKanbanChainHandler: Finished successfully")
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

		var kanbanChainRequest struct {
			KanbanChain        models.KanbanChain `json:"kanban_chain"`
			NoOfInitialKanbans int64              `json:"no_of_initial_kanbans"`
		}

		err = json.NewDecoder(r.Body).Decode(&kanbanChainRequest) // Decode the data once
		if err != nil {
			log.Printf("UpdateKanbanChainHandler: Error reading request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		log.Printf("UpdateKanbanChainHandler: Request body decoded, Data: %+v", kanbanChainRequest)

		kanbanChainUpdates := kanbanChainRequest.KanbanChain
		kanbanChainUpdates.ID = id // Ensure ID from URL is used
		log.Printf("UpdateKanbanChainHandler: Updating kanban chain: %+v", kanbanChainUpdates)

		updatedKanbanChain, err := updateKanbanChain(db, kanbanChainUpdates)
		if err != nil {
			log.Printf("UpdateKanbanChainHandler: Failed to update kanban chain: %v", err)
			http.Error(w, "Failed to update kanban chain", http.StatusInternalServerError)
			return
		}

		log.Printf("UpdateKanbanChainHandler: Kanban chain updated successfully with data: %+v", updatedKanbanChain)

		if kanbanChainRequest.NoOfInitialKanbans > 0 {
			err = createInitialKanbans(db, updatedKanbanChain.ID, kanbanChainRequest.NoOfInitialKanbans, updatedKanbanChain.StatusChainID, updatedKanbanChain.LeadtimeDays, updatedKanbanChain.TipoContenitore, updatedKanbanChain.Quantity)
			if err != nil {
				log.Printf("UpdateKanbanChainHandler: Error creating additional kanbans: %v", err)
				http.Error(w, "Failed to create initial kanbans for the chain", http.StatusInternalServerError)
				return
			}
			log.Println("UpdateKanbanChainHandler: initial kanbans created successfully")

			// Retrieve the updated kanban chain to get the correct no_of_active_kanbans count
			updatedKanbanChain, err = getKanbanChainByID(db, id)
			if err != nil {
				log.Printf("UpdateKanbanChainHandler: Error re-fetching kanban chain after kanban creation: %v", err)
				// Not critical, but log the error
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedKanbanChain)
		log.Println("UpdateKanbanChainHandler: Finished successfully")

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

func getKanbanChains(db *sql.DB) ([]map[string]interface{}, error) {
	rows, err := db.Query(`
		SELECT
			kc.id,
			c.name AS customer_name,
			kc.cliente_id,
			p.name AS product_name,
            kc.prodotto_codice,
			s.name AS supplier_name,
			kc.fornitore_id,
			kc.leadtime_days,
			kc.quantity,
			kc.tipo_contenitore,
			kc.status_chain_id,
			kc.no_of_active_kanbans
		FROM kanban_chains kc
		JOIN accounts c ON kc.cliente_id = c.id
		JOIN products p ON kc.prodotto_codice = p.product_id
        JOIN accounts s ON kc.fornitore_id = s.id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var kanbanChains []map[string]interface{}
	for rows.Next() {
		var kcID int64
		var customerName string
		var clienteID int64
		var productName string
		var prodottoCodice string
		var supplierName string
		var fornitoreID int64
		var leadtimeDays int64
		var quantity float64
		var tipoContenitore string
		var statusChainID int64
		var noOfActiveKanbans int64

		if err := rows.Scan(
			&kcID, &customerName, &clienteID, &productName, &prodottoCodice, &supplierName, &fornitoreID,
			&leadtimeDays, &quantity, &tipoContenitore, &statusChainID, &noOfActiveKanbans,
		); err != nil {
			return nil, err
		}
		kanbanChains = append(kanbanChains, map[string]interface{}{
			"id":                   kcID,
			"customer_name":        customerName,
			"cliente_id":           clienteID,
			"product_name":         productName,
			"prodotto_codice":      prodottoCodice,
			"supplier_name":        supplierName,
			"fornitore_id":         fornitoreID,
			"leadtime_days":        leadtimeDays,
			"quantity":             quantity,
			"tipo_contenitore":     tipoContenitore,
			"status_chain_id":      statusChainID,
			"no_of_active_kanbans": noOfActiveKanbans,
		})
	}
	return kanbanChains, nil
}

func createKanbanChain(db *sql.DB, kc models.KanbanChain) (*models.KanbanChain, error) {
	log.Printf("createKanbanChain: Starting with data: %+v", kc)
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
	log.Printf("createKanbanChain: SQL Query: %s, Parameters: [%d, %s, %d, %d, %f, %s, %d, %d]", sqlStatement, kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays, kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans)
	var newKC models.KanbanChain
	err := db.QueryRow(sqlStatement,
		kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays,
		kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans,
	).Scan(
		&newKC.ID, &newKC.ClienteID, &newKC.ProdottoCodice, &newKC.FornitoreID, &newKC.LeadtimeDays,
		&newKC.Quantity, &newKC.TipoContenitore, &newKC.StatusChainID, &newKC.NoOfActiveKanbans,
	)
	if err != nil {
		log.Printf("createKanbanChain: Error executing query: %v", err)
		return nil, err
	}
	log.Printf("createKanbanChain: Kanban chain created successfully with data: %+v", newKC)

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
	log.Printf("updateKanbanChain: Starting with data: %+v", kc)
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
	log.Printf("updateKanbanChain: SQL Query: %s, Parameters: [%d, %d, %s, %d, %d, %f, %s, %d, %d]", sqlStatement, kc.ID, kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays, kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans)
	var updatedKC models.KanbanChain
	err := db.QueryRow(sqlStatement,
		kc.ID, kc.ClienteID, kc.ProdottoCodice, kc.FornitoreID, kc.LeadtimeDays,
		kc.Quantity, kc.TipoContenitore, kc.StatusChainID, kc.NoOfActiveKanbans,
	).Scan(
		&updatedKC.ID, &updatedKC.ClienteID, &updatedKC.ProdottoCodice, &updatedKC.FornitoreID, &updatedKC.LeadtimeDays,
		&updatedKC.Quantity, &updatedKC.TipoContenitore, &updatedKC.StatusChainID, &updatedKC.NoOfActiveKanbans,
	)
	if err != nil {
		log.Printf("updateKanbanChain: Error executing query: %v", err)
		return nil, err
	}
	log.Printf("updateKanbanChain: Kanban chain updated successfully with data: %+v", updatedKC)
	return &updatedKC, nil
}

func deleteKanbanChain(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM kanban_chains WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}

// createInitialKanbans creates kanban records when a new kanban chain is created
func createInitialKanbans(db *sql.DB, kanbanChainID int64, numberOfKanbans int64, statusChainID int64, leadtimeDays int64, tipoContenitore string, quantity float64) error {
	log.Println("createInitialKanbans: Starting") // Log entry
	log.Printf("createInitialKanbans: Parameters - kanbanChainID: %d, numberOfKanbans: %d, statusChainID: %d, leadtimeDays: %d, tipoContenitore: %s, quantity: %f", kanbanChainID, numberOfKanbans, statusChainID, leadtimeDays, tipoContenitore, quantity)

	tx, err := db.Begin()
	if err != nil {
		log.Printf("createInitialKanbans: Error starting transaction: %v", err) // Log transaction begin error
		return fmt.Errorf("error starting transaction for creating initial kanbans: %w", err)
	}
	defer tx.Rollback() // Rollback if any operation fails

	// Get the first status in the status chain to set as initial status_current
	firstStatusID, err := getFirstStatusIDInChain(db, statusChainID)
	if err != nil {
		log.Printf("createInitialKanbans: Error getting first status in chain: %v", err) // Log error getting first status
		return fmt.Errorf("error getting first status in chain: %w", err)
	}
	if firstStatusID == 0 {
		log.Printf("createInitialKanbans: No statuses found in status chain %d", statusChainID) // Log no statuses found
		return fmt.Errorf("no statuses found in status chain %d", statusChainID)
	}
	log.Printf("createInitialKanbans: First Status ID in chain %d: %d", statusChainID, firstStatusID) // Log first status ID

	sqlStatement := `
		INSERT INTO kanbans (
			kanban_chain_id, status_chain_id, status_current, leadtime_days, tipo_contenitore, quantity, data_aggiornamento
		)
		SELECT $1, $2, $3, $4, $5, $6, NOW()
		FROM generate_series(1, $7)
	` // data_aggiornamento set to NOW() on creation
	log.Printf("createInitialKanbans: SQL Query: %s, Parameters: [kanbanChainID=%d, statusChainID=%d, firstStatusID=%d, leadtimeDays=%d, tipoContenitore=%s, quantity=%f, numberOfKanbans=%d]", sqlStatement, kanbanChainID, statusChainID, firstStatusID, leadtimeDays, tipoContenitore, quantity, numberOfKanbans)

	_, err = tx.Exec(sqlStatement, kanbanChainID, statusChainID, firstStatusID, leadtimeDays, tipoContenitore, quantity, numberOfKanbans)
	if err != nil {
		log.Printf("createInitialKanbans: Error executing kanban insert query: %v", err) // Log query execution error
		return fmt.Errorf("error inserting kanbans: %w", err)
	}
	log.Println("createInitialKanbans: Kanban insert query executed successfully") // Log query success

	if err := tx.Commit(); err != nil {
		log.Printf("createInitialKanbans: Error committing transaction: %v", err) // Log transaction commit error
		return fmt.Errorf("error committing transaction for initial kanban creation: %w", err)
	}

	log.Println("createInitialKanbans: Finished successfully") // Log exit
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
