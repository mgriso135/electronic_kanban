package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"electronic_kanban_backend/models"

	"github.com/gorilla/mux"
)

// GetProductsHandler returns a handler for GET /api/products
func GetProductsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		products, err := getProducts(db)
		if err != nil {
			http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(products)
	}
}

// CreateProductHandler returns a handler for POST /api/products
func CreateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var product models.Product
		err := json.NewDecoder(r.Body).Decode(&product)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		newProduct, err := createProduct(db, product)
		if err != nil {
			http.Error(w, "Failed to create product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newProduct)
	}
}

// GetProductHandler returns a handler for GET /api/products/{id}
func GetProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // Product ID is text, no need to parse to int

		product, err := getProductByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.NotFound(w, r)
			} else {
				http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(product)
	}
}

// UpdateProductHandler returns a handler for PUT/PATCH /api/products/{id}
func UpdateProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // Product ID is text

		var productUpdates models.Product
		err := json.NewDecoder(r.Body).Decode(&productUpdates)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		productUpdates.ProductID = id // Ensure ID from URL is used

		updatedProduct, err := updateProduct(db, productUpdates)
		if err != nil {
			http.Error(w, "Failed to update product", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(updatedProduct)
	}
}

// DeleteProductHandler returns a handler for DELETE /api/products/{id}
func DeleteProductHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id := vars["id"] // Product ID is text

		err := deleteProduct(db, id)
		if err != nil {
			http.Error(w, "Failed to delete product", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK) // Respond with 200 OK
		json.NewEncoder(w).Encode(map[string]string{"message": "Product deleted"})
	}
}

// Database interaction functions (private)

func getProducts(db *sql.DB) ([]models.Product, error) {
	rows, err := db.Query("SELECT product_id, name FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var product models.Product
		if err := rows.Scan(&product.ProductID, &product.Name); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, nil
}

func createProduct(db *sql.DB, product models.Product) (*models.Product, error) {
	sqlStatement := `
		INSERT INTO products (product_id, name)
		VALUES ($1, $2)
		RETURNING product_id, name`
	var newProduct models.Product
	err := db.QueryRow(sqlStatement, product.ProductID, product.Name).Scan(
		&newProduct.ProductID, &newProduct.Name,
	)
	if err != nil {
		return nil, err
	}
	return &newProduct, nil
}

func getProductByID(db *sql.DB, id string) (*models.Product, error) {
	sqlStatement := `SELECT product_id, name FROM products WHERE product_id = $1`
	var product models.Product
	err := db.QueryRow(sqlStatement, id).Scan(&product.ProductID, &product.Name)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func updateProduct(db *sql.DB, product models.Product) (*models.Product, error) {
	sqlStatement := `
		UPDATE products
		SET name = $2
		WHERE product_id = $1
		RETURNING product_id, name`
	var updatedProduct models.Product
	err := db.QueryRow(sqlStatement, product.ProductID, product.Name).Scan(
		&updatedProduct.ProductID, &updatedProduct.Name,
	)
	if err != nil {
		return nil, err
	}
	return &updatedProduct, nil
}

func deleteProduct(db *sql.DB, id string) error {
	sqlStatement := `DELETE FROM products WHERE product_id = $1`
	_, err := db.Exec(sqlStatement, id)
	return err
}
