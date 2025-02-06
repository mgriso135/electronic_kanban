package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"electronic_kanban_backend/db"
	"electronic_kanban_backend/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors" // Import the cors package
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file, using environment variables if set")
	}

	// Get database connection string from the environment or use the default value
	dbConnStr := os.Getenv("DB_CONNECTION_STRING")
	if dbConnStr == "" {
		dbConnStr = "postgres://postgres:password@localhost:5432/electronic_kanban?sslmode=disable"
		log.Println("DB_CONNECTION_STRING not set, using default for local development (postgres://postgres:password@localhost:5432/electronic_kanban?sslmode=disable)")
	}
	fmt.Printf("Database connection string is: %s\n", dbConnStr) // DEBUG: Print the connection string

	log.Println("Attempting to connect to the database...") //DEBUG: Before db connection
	// Initialize database connection
	database, err := db.ConnectDB(dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()
	log.Println("Successfully connected to the database") //DEBUG: After db connection

	router := mux.NewRouter()

	// Account Routes
	router.HandleFunc("/api/accounts", handlers.GetAccountsHandler(database)).Methods("GET")
	router.HandleFunc("/api/accounts", handlers.CreateAccountHandler(database)).Methods("POST")
	router.HandleFunc("/api/accounts/{id}", handlers.GetAccountHandler(database)).Methods("GET")
	router.HandleFunc("/api/accounts/{id}", handlers.UpdateAccountHandler(database)).Methods("PUT", "PATCH")
	router.HandleFunc("/api/accounts/{id}", handlers.DeleteAccountHandler(database)).Methods("DELETE")

	// Product Routes
	router.HandleFunc("/api/products", handlers.GetProductsHandler(database)).Methods("GET")
	router.HandleFunc("/api/products", handlers.CreateProductHandler(database)).Methods("POST")
	router.HandleFunc("/api/products/{id}", handlers.GetProductHandler(database)).Methods("GET")
	router.HandleFunc("/api/products/{id}", handlers.UpdateProductHandler(database)).Methods("PUT", "PATCH")
	router.HandleFunc("/api/products/{id}", handlers.DeleteProductHandler(database)).Methods("DELETE")

	// Status Routes
	router.HandleFunc("/api/statuses", handlers.GetStatusesHandler(database)).Methods("GET")
	router.HandleFunc("/api/statuses", handlers.CreateStatusHandler(database)).Methods("POST")
	router.HandleFunc("/api/statuses/{id}", handlers.GetStatusHandler(database)).Methods("GET")
	router.HandleFunc("/api/statuses/{id}", handlers.UpdateStatusHandler(database)).Methods("PUT", "PATCH")
	router.HandleFunc("/api/statuses/{id}", handlers.DeleteStatusHandler(database)).Methods("DELETE")

	// Status Chain Routes
	router.HandleFunc("/api/status-chains", handlers.GetStatusChainsHandler(database)).Methods("GET")
	router.HandleFunc("/api/status-chains", handlers.CreateStatusChainHandler(database)).Methods("POST")
	router.HandleFunc("/api/status-chains/{id}", handlers.GetStatusChainHandler(database)).Methods("GET")
	router.HandleFunc("/api/status-chains/{id}", handlers.UpdateStatusChainHandler(database)).Methods("PUT", "PATCH")
	router.HandleFunc("/api/status-chains/{id}", handlers.DeleteStatusChainHandler(database)).Methods("DELETE")
	router.HandleFunc("/api/status-chains/{statusChainId}/statuses", handlers.GetStatusChainStatusesHandler(database)).Methods("GET")
	router.HandleFunc("/api/status-chains/{statusChainId}/statuses", handlers.UpdateStatusChainStatusesHandler(database)).Methods("PUT")

	// Kanban Chain Routes
	router.HandleFunc("/api/kanban-chains", handlers.GetKanbanChainsHandler(database)).Methods("GET")
	router.HandleFunc("/api/kanban-chains", handlers.CreateKanbanChainHandler(database)).Methods("POST")
	router.HandleFunc("/api/kanban-chains/{id}", handlers.GetKanbanChainHandler(database)).Methods("GET")
	router.HandleFunc("/api/kanban-chains/{id}", handlers.UpdateKanbanChainHandler(database)).Methods("PUT", "PATCH")
	router.HandleFunc("/api/kanban-chains/{id}", handlers.DeleteKanbanChainHandler(database)).Methods("DELETE")

	router.HandleFunc("/api/kanbans", handlers.GetKanbansHandler(database)).Methods("GET") // GET with optional product filter
	router.HandleFunc("/api/kanbans", handlers.CreateKanbanHandler(database)).Methods("POST")
	router.HandleFunc("/api/kanbans/{id}", handlers.GetKanbanHandler(database)).Methods("GET")             // GetKanbanHandler for GET
	router.HandleFunc("/api/kanbans/{id}", handlers.UpdateKanbanHandler(database)).Methods("PUT", "PATCH") // UpdateKanbanHandler for PUT/PATCH
	router.HandleFunc("/api/kanbans/{id}", handlers.DeleteKanbanHandler(database)).Methods("DELETE")

	// Dashboard Routes
	router.HandleFunc("/api/dashboards/supplier/{supplierId}", handlers.GetSupplierDashboardHandler(database)).Methods("GET")
	router.HandleFunc("/api/dashboards/customer/{customerId}", handlers.GetCustomerDashboardHandler(database)).Methods("GET")

	// Setup CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow requests from your React frontend
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"}, // or you can specify specific headers if needed
	})

	// Apply the CORS middleware to all routes
	handler := c.Handler(router)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port if not specified in environment
	}
	fmt.Printf("Server listening on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler)) // use handler here
}
