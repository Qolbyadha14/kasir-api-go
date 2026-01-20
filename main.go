package main

import (
	"encoding/json"
	"fmt"
	"net/http"
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

func main() {
	// Get localhost:8080/health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		// Response JSON
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"message": "API Running",
		})
	})

	// Handle /api/products (GET and POST)
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Starting server on port 8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
