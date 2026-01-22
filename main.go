package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api-go/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

type Product struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
	Stock int    `json:"stock"`
}

var products = []Product{
	{
		ID:    1,
		Name:  "Product 1",
		Price: 10000,
		Stock: 10,
	},
	{
		ID:    2,
		Name:  "Product 2",
		Price: 20000,
		Stock: 20,
	},
}

func getProduct(id int) (Product, bool) {
	for _, p := range products {
		if p.ID == id {
			return p, true
		}
	}
	return Product{}, false
}

func updateProduct(id int, product Product) bool {
	for i, p := range products {
		if p.ID == id {
			products[i] = product
			return true
		}
	}
	return false
}

func deleteProduct(id int) bool {
	for i, p := range products {
		if p.ID == id {
			products = append(products[:i], products[i+1:]...)
			return true
		}
	}
	return false
}

// @Summary Health check
// @Description Get the status of the API
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	// Response JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":  "ok",
		"message": "API Running",
	})
}

// @Summary Create or List products
// @Description Handle /api/products (GET and POST)
// @Tags products
// @Accept json
// @Produce json
// @Param product body Product false "Product object (for POST)"
// @Success 201 {object} Product
// @Success 200 {array} Product
// @Failure 400 {object} map[string]string
// @Router /api/products [post]
// @Router /api/products [get]
func ProductsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Method Post
	if r.Method == http.MethodPost {
		var product Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		products = append(products, product)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(product)
		return
	}

	json.NewEncoder(w).Encode(products)
}

// @Summary Get, Update or Delete a product
// @Description Handle /api/products/{id} (GET, UPDATE AND DELETE)
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body Product false "Product object (for PUT)"
// @Success 200 {object} Product
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /api/products/{id} [get]
// @Router /api/products/{id} [put]
// @Router /api/products/{id} [delete]
func ProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	if idStr == "" {
		return
	}

	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "invalid product id"})
		return
	}

	// Method Get
	if r.Method == http.MethodGet {
		if product, found := getProduct(id); found {
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	// Method Put
	if r.Method == http.MethodPut {
		var product Product
		if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "invalid request body"})
			return
		}

		product.ID = id
		if ok := updateProduct(id, product); ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(product)
			return
		}
	}

	// Method Delete
	if r.Method == http.MethodDelete {
		if ok := deleteProduct(id); ok {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{"message": "product deleted"})
			return
		}
	}

	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(map[string]string{"error": "product not found"})
}

// @title Kasir API
// @version 1.0
// @description This is a simple Kasir API server.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {
	// Get localhost:8080/health
	http.HandleFunc("/health", HealthHandler)

	// Handle /api/products (GET and POST)
	http.HandleFunc("/api/products", ProductsHandler)

	// Handle /api/products/{id} (GET, UPDATE AND DELETE)
	http.HandleFunc("/api/products/", ProductDetailHandler)

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Swagger UI
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
