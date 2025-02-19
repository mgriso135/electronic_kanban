package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// GetSupplierDashboardHandler will handle GET requests to /api/dashboards/supplier/{supplierId}
func GetSupplierDashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		supplierIDStr := vars["supplierId"]
		supplierID, err := strconv.ParseInt(supplierIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid supplier ID", http.StatusBadRequest)
			return
		}

		kanbans, err := getKanbansForSupplierDashboard(db, supplierID)
		if err != nil {
			http.Error(w, "Failed to fetch kanban data for supplier dashboard", http.StatusInternalServerError)
			return
		}

		// Organize kanbans by product (you might want to refine this structure)
		kanbansByProduct := make(map[string][]map[string]interface{})
		for _, kanban := range kanbans {
			productID, ok := kanban["prodotto_codice"].(string)
			if !ok {
				productID = "Unknown Product" // Handle case where product_id is not a string
			}
			kanbansByProduct[productID] = append(kanbansByProduct[productID], kanban)
		}

		response := map[string]interface{}{
			"supplier_id":        supplierID,
			"kanbans_by_product": kanbansByProduct,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// GetCustomerDashboardHandler will handle GET requests to /api/dashboards/customer/{customerId}
func GetCustomerDashboardHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		customerIDStr := vars["customerId"]
		customerID, err := strconv.ParseInt(customerIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Invalid customer ID", http.StatusBadRequest)
			return
		}

		kanbans, err := getKanbansForCustomerDashboard(db, customerID)
		if err != nil {
			http.Error(w, "Failed to fetch kanban data for customer dashboard", http.StatusInternalServerError)
			return
		}

		// Organize kanbans by product (you might want to refine this structure)
		kanbansByProduct := make(map[string][]map[string]interface{})
		for _, kanban := range kanbans {
			productID, ok := kanban["prodotto_codice"].(string)
			if !ok {
				productID = "Unknown Product" // Handle case where product_id is not a string
			}
			kanbansByProduct[productID] = append(kanbansByProduct[productID], kanban)
		}

		response := map[string]interface{}{
			"customer_id":        customerID,
			"kanbans_by_product": kanbansByProduct,
		}

		json.NewEncoder(w).Encode(response)
	}
}

// Database interaction functions (private)

// getKanbansForSupplierDashboard retrieves Kanban data for the supplier dashboard.
func getKanbansForSupplierDashboard(db *sql.DB, supplierID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT
			k.id AS kanban_id,
			kc.prodotto_codice,
			p.name AS product_name,
			k.tipo_contenitore,
			k.quantity,
			s.name AS status_name,
			s.color AS status_color,
			scs.customer_supplier,
            k.status_current,
			ac.name AS customer_name  -- ADD CUSTOMER NAME HERE - as per user request, but query was for supplier dashboard, so showing customer name here
		FROM
			kanbans k
		JOIN
			kanban_chains kc ON k.kanban_chain_id = kc.id
		JOIN
			products p ON kc.prodotto_codice = p.product_id
		JOIN
			statuses s ON k.status_current = s.status_id
		JOIN
			status_chains_statuses scs ON k.status_chain_id = scs.status_chain_id AND k.status_current = scs.status_id
		JOIN
			accounts ac ON kc.cliente_id = ac.id  -- JOIN with accounts table for customer name -  corrected JOIN for customer name for supplier dashboard
		WHERE
			kc.fornitore_id = $1  -- WHERE clause for supplier dashboard
		ORDER BY
			p.name, k.id;
	`

	rows, err := db.Query(query, supplierID)
	if err != nil {
		return nil, fmt.Errorf("error querying kanbans for supplier dashboard: %w", err)
	}
	defer rows.Close()

	var kanbans []map[string]interface{}
	for rows.Next() {
		var kanbanID int64
		var prodottoCodice string
		var productName string
		var tipoContenitore string
		var quantity float64
		var statusName string
		var statusColor string
		var customerSupplier int
		var statusCurrent int64
		var customerName string // Variable for customer name

		if err := rows.Scan(
			&kanbanID, &prodottoCodice, &productName, &tipoContenitore, &quantity, &statusName, &statusColor, &customerSupplier,
			&statusCurrent, &customerName, // Scan customerName here
		); err != nil {
			return nil, fmt.Errorf("error scanning kanban row for supplier dashboard: %w", err)
		}

		kanbans = append(kanbans, map[string]interface{}{
			"kanban_id":         kanbanID,
			"prodotto_codice":   prodottoCodice,
			"product_name":      productName,
			"tipo_contenitore":  tipoContenitore,
			"quantity":          quantity,
			"status_name":       statusName,
			"status_color":      statusColor,
			"customer_supplier": customerSupplier,
			"status_current":    statusCurrent,
			"customer_name":     customerName, // Add customer_name to the map
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kanban rows for supplier dashboard: %w", err)
	}

	return kanbans, nil
}

// getKanbansForCustomerDashboard retrieves Kanban data for the customer dashboard.
func getKanbansForCustomerDashboard(db *sql.DB, customerID int64) ([]map[string]interface{}, error) {
	query := `
		SELECT
			k.id AS kanban_id,
			kc.prodotto_codice,
			p.name AS product_name,
			k.tipo_contenitore,
			k.quantity,
			s.name AS status_name,
			s.color AS status_color,
			scs.customer_supplier,
            k.status_current,
			ac.name AS supplier_name  -- ADD SUPPLIER NAME HERE
		FROM
			kanbans k
		JOIN
			kanban_chains kc ON k.kanban_chain_id = kc.id
		JOIN
			products p ON kc.prodotto_codice = p.product_id
		JOIN
			statuses s ON k.status_current = s.status_id
		JOIN
			status_chains_statuses scs ON k.status_chain_id = scs.status_chain_id AND k.status_current = scs.status_id
		JOIN
			accounts ac ON kc.fornitore_id = ac.id  -- JOIN with accounts table for supplier name
		WHERE
			kc.cliente_id = $1
		ORDER BY
			p.name, k.id;
	`

	rows, err := db.Query(query, customerID)
	if err != nil {
		return nil, fmt.Errorf("error querying kanbans for customer dashboard: %w", err)
	}
	defer rows.Close()

	var kanbans []map[string]interface{}
	for rows.Next() {
		var kanbanID int64
		var prodottoCodice string
		var productName string
		var tipoContenitore string
		var quantity float64
		var statusName string
		var statusColor string
		var customerSupplier int
		var statusCurrent int64
		var supplierName string // Variable for supplier name

		if err := rows.Scan(
			&kanbanID, &prodottoCodice, &productName, &tipoContenitore, &quantity, &statusName, &statusColor, &customerSupplier,
			&statusCurrent, &supplierName, // Scan supplierName here
		); err != nil {
			return nil, fmt.Errorf("error scanning kanban row for customer dashboard: %w", err)
		}

		kanbans = append(kanbans, map[string]interface{}{
			"kanban_id":         kanbanID,
			"prodotto_codice":   prodottoCodice,
			"product_name":      productName,
			"tipo_contenitore":  tipoContenitore,
			"quantity":          quantity,
			"status_name":       statusName,
			"status_color":      statusColor,
			"customer_supplier": customerSupplier,
			"status_current":    statusCurrent,
			"supplier_name":     supplierName, // Add supplier_name to the map
		})
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating kanban rows for customer dashboard: %w", err)
	}

	return kanbans, nil
}
