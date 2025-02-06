package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"electronic_kanban_backend/models"

	"github.com/gorilla/mux"
)

// GetAccountsHandler returns a handler for GET /api/accounts
func GetAccountsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetAccountsHandler: Starting")
		accounts, err := getAccounts(db)
		if err != nil {
			log.Printf("GetAccountsHandler: Error fetching accounts: %v", err)
			http.Error(w, "Failed to fetch accounts", http.StatusInternalServerError)
			return
		}
		log.Println("GetAccountsHandler: Finished successfully")
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(accounts)
	}
}

// CreateAccountHandler returns a handler for POST /api/accounts
func CreateAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("CreateAccountHandler: Starting")
		var account models.Account
		err := json.NewDecoder(r.Body).Decode(&account)
		if err != nil {
			log.Printf("CreateAccountHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		log.Printf("CreateAccountHandler: Creating account with data: %+v", account)

		newAccount, err := createAccount(db, account)
		if err != nil {
			log.Printf("CreateAccountHandler: Failed to create account: %v", err)
			http.Error(w, "Failed to create account", http.StatusInternalServerError)
			return
		}
		log.Printf("CreateAccountHandler: Account created successfully with data: %+v", newAccount)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newAccount)
	}
}

// GetAccountHandler returns a handler for GET /api/accounts/{id}
func GetAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("GetAccountHandler: Starting")
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Println("GetAccountHandler: Account ID is required")
			http.Error(w, "Account ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("GetAccountHandler: Invalid account ID: %v", err)
			http.Error(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		log.Printf("GetAccountHandler: Getting account with ID: %d", id)

		account, err := getAccountByID(db, id)
		if err != nil {
			log.Printf("GetAccountHandler: Error fetching account with ID %d: %v", id, err)
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch account", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("GetAccountHandler: Successfully retrieved account: %+v", account)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(account)
	}
}

// UpdateAccountHandler returns a handler for PUT/PATCH /api/accounts/{id}
func UpdateAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateAccountHandler: Starting")
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Println("UpdateAccountHandler: Account ID is required")
			http.Error(w, "Account ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("UpdateAccountHandler: Invalid account ID: %v", err)
			http.Error(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		var accountUpdates models.Account
		err = json.NewDecoder(r.Body).Decode(&accountUpdates)
		if err != nil {
			log.Printf("UpdateAccountHandler: Invalid request body: %v", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		accountUpdates.ID = id // Ensure ID from URL is used
		log.Printf("UpdateAccountHandler: Updating account with data: %+v", accountUpdates)

		updatedAccount, err := updateAccount(db, accountUpdates)
		if err != nil {
			log.Printf("UpdateAccountHandler: Failed to update account: %v", err)
			http.Error(w, "Failed to update account", http.StatusInternalServerError)
			return
		}
		log.Printf("UpdateAccountHandler: Successfully updated account: %+v", updatedAccount)

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedAccount)
	}
}

// DeleteAccountHandler returns a handler for DELETE /api/accounts/{id}
func DeleteAccountHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("DeleteAccountHandler: Starting")
		vars := mux.Vars(r)
		idStr, ok := vars["id"]
		if !ok {
			log.Println("DeleteAccountHandler: Account ID is required")
			http.Error(w, "Account ID is required", http.StatusBadRequest)
			return
		}
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			log.Printf("DeleteAccountHandler: Invalid account ID: %v", err)
			http.Error(w, "Invalid account ID", http.StatusBadRequest)
			return
		}

		log.Printf("DeleteAccountHandler: Deleting account with ID: %d", id)

		err = deleteAccount(db, id)
		if err != nil {
			log.Printf("DeleteAccountHandler: Failed to delete account: %v", err)
			http.Error(w, "Failed to delete account", http.StatusInternalServerError)
			return
		}
		log.Println("DeleteAccountHandler: Account deleted successfully")

		w.WriteHeader(http.StatusOK) // Respond with 200 OK for successful deletion
		json.NewEncoder(w).Encode(map[string]string{"message": "Account deleted"})
	}
}

// Database interaction functions (private)

func getAccounts(db *sql.DB) ([]models.Account, error) {
	log.Println("getAccounts: Starting")
	rows, err := db.Query("SELECT id, name, vat_number, address FROM accounts")
	if err != nil {
		log.Printf("getAccounts: Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var accounts []models.Account
	for rows.Next() {
		var account models.Account
		if err := rows.Scan(&account.ID, &account.Name, &account.VATNumber, &account.Address); err != nil {
			log.Printf("getAccounts: Error scanning rows: %v", err)
			return nil, err
		}
		accounts = append(accounts, account)
	}
	log.Printf("getAccounts: Successfully retrieved %d accounts", len(accounts))

	return accounts, nil
}

func createAccount(db *sql.DB, account models.Account) (*models.Account, error) {
	log.Printf("createAccount: Starting with account data: %+v", account)
	sqlStatement := `
		INSERT INTO accounts (name, vat_number, address)
		VALUES ($1, $2, $3)
		RETURNING id, name, vat_number, address`
	var newAccount models.Account
	err := db.QueryRow(sqlStatement, account.Name, account.VATNumber, account.Address).Scan(
		&newAccount.ID, &newAccount.Name, &newAccount.VATNumber, &newAccount.Address,
	)
	if err != nil {
		log.Printf("createAccount: Error creating account: %v", err)
		return nil, err
	}
	log.Printf("createAccount: Successfully created account with data: %+v", newAccount)

	return &newAccount, nil
}

func getAccountByID(db *sql.DB, id int64) (*models.Account, error) {
	log.Printf("getAccountByID: Starting with ID: %d", id)
	sqlStatement := `SELECT id, name, vat_number, address FROM accounts WHERE id = $1`
	var account models.Account
	err := db.QueryRow(sqlStatement, id).Scan(&account.ID, &account.Name, &account.VATNumber, &account.Address)
	if err != nil {
		log.Printf("getAccountByID: Error querying database with ID %d: %v", id, err)
		return nil, err
	}
	log.Printf("getAccountByID: Successfully retrieved account: %+v", account)
	return &account, nil
}

func updateAccount(db *sql.DB, account models.Account) (*models.Account, error) {
	log.Printf("updateAccount: Starting with account data: %+v", account)
	sqlStatement := `
		UPDATE accounts
		SET name = $2, vat_number = $3, address = $4
		WHERE id = $1
		RETURNING id, name, vat_number, address`
	var updatedAccount models.Account
	err := db.QueryRow(sqlStatement, account.ID, account.Name, account.VATNumber, account.Address).Scan(
		&updatedAccount.ID, &updatedAccount.Name, &updatedAccount.VATNumber, &updatedAccount.Address,
	)
	if err != nil {
		log.Printf("updateAccount: Error updating account: %v", err)
		return nil, err
	}
	log.Printf("updateAccount: Successfully updated account with data: %+v", updatedAccount)
	return &updatedAccount, nil
}

func deleteAccount(db *sql.DB, id int64) error {
	log.Printf("deleteAccount: Starting with ID: %d", id)
	sqlStatement := `DELETE FROM accounts WHERE id = $1`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Printf("deleteAccount: Error deleting account: %v", err)
		return err
	}

	log.Printf("deleteAccount: Successfully deleted account")

	return nil
}
