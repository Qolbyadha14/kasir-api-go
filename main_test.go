package main

import (
	"bytes"
	"encoding/json"
	"kasir-api-go/internal/api"
	"kasir-api-go/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Helper to decode response
func decodeResponse(t *testing.T, body *bytes.Buffer) api.JSONResponse {
	var response api.JSONResponse
	err := json.NewDecoder(body).Decode(&response)
	if err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	return response
}

func TestGetProduct(t *testing.T) {
	// Unit tests for helper functions can remain the same
	// Test Case 1: Product found
	t.Run("Found", func(t *testing.T) {
		id := 1
		product, found := getProduct(id)

		if !found {
			t.Errorf("Expected product with ID %d to be found", id)
		}

		if product.ID != id {
			t.Errorf("Expected product ID %d, got %d", id, product.ID)
		}

		if product.Name == "" {
			t.Error("Expected product name to be not empty")
		}
	})

	// Test Case 2: Product not found
	t.Run("NotFound", func(t *testing.T) {
		id := 999
		_, found := getProduct(id)

		if found {
			t.Errorf("Expected product with ID %d to be not found", id)
		}
	})
}

// ... (Other helper function tests remain unchanged) ...
func TestGetCategory(t *testing.T) {
	// Test Case 1: Category found
	t.Run("Found", func(t *testing.T) {
		id := 1
		category, found := getCategory(id)

		if !found {
			t.Errorf("Expected category with ID %d to be found", id)
		}

		if category.ID != id {
			t.Errorf("Expected category ID %d, got %d", id, category.ID)
		}
	})

	// Test Case 2: Category not found
	t.Run("NotFound", func(t *testing.T) {
		id := 999
		_, found := getCategory(id)

		if found {
			t.Errorf("Expected category with ID %d to be not found", id)
		}
	})
}

func TestUpdateProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1
		originalProduct, _ := getProduct(id)
		updatedProduct := models.Product{
			ID:    id,
			Name:  "Updated Product 1",
			Price: 15000,
			Stock: 5,
		}

		ok := updateProduct(id, updatedProduct)
		if !ok {
			t.Errorf("Expected updateProduct to return true for ID %d", id)
		}

		p, _ := getProduct(id)
		if p.Name != "Updated Product 1" {
			t.Errorf("Expected name 'Updated Product 1', got %s", p.Name)
		}

		// Restore original for other tests
		updateProduct(id, originalProduct)
	})

	t.Run("NotFound", func(t *testing.T) {
		id := 999
		ok := updateProduct(id, models.Product{ID: id, Name: "Non-existent"})
		if ok {
			t.Error("Expected updateProduct to return false for non-existent ID")
		}
	})
}

func TestDeleteProduct(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		// Add a temporary product to delete
		tempProduct := models.Product{ID: 100, Name: "Delete Me"}
		products = append(products, tempProduct)

		ok := deleteProduct(100)
		if !ok {
			t.Error("Expected deleteProduct to return true")
		}

		_, found := getProduct(100)
		if found {
			t.Error("Expected product to be deleted from slice")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ok := deleteProduct(999)
		if ok {
			t.Error("Expected deleteProduct to return false for non-existent ID")
		}
	})
}

func TestUpdateCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		id := 1
		originalCategory, _ := getCategory(id)
		updatedCategory := models.Category{
			ID:          id,
			Name:        "Updated Category 1",
			Description: "Updated Description",
		}

		ok := updateCategory(id, updatedCategory)
		if !ok {
			t.Errorf("Expected updateCategory to return true for ID %d", id)
		}

		c, _ := getCategory(id)
		if c.Name != "Updated Category 1" {
			t.Errorf("Expected name 'Updated Category 1', got %s", c.Name)
		}

		// Restore original
		updateCategory(id, originalCategory)
	})
}

func TestDeleteCategory(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		tempCategory := models.Category{ID: 100, Name: "Delete Category"}
		categories = append(categories, tempCategory)

		ok := deleteCategory(100)
		if !ok {
			t.Error("Expected deleteCategory to return true")
		}

		_, found := getCategory(100)
		if found {
			t.Error("Expected category to be deleted from slice")
		}
	})
}

// Updated Handler Tests
func TestCreateProductHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		newProduct := models.Product{
			ID:    10,
			Name:  "Product 10",
			Price: 50000,
			Stock: 50,
		}
		body, _ := json.Marshal(newProduct)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(CreateProductHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		response := decodeResponse(t, rr.Body)
		if !response.Success {
			t.Error("Expected success to be true")
		}

		// We need to re-marshal keys to map to struct for verification if needed
		// Or simpler, just check success flag which is enough for now
	})

	t.Run("DuplicateID", func(t *testing.T) {
		existingProduct := products[0]
		body, _ := json.Marshal(existingProduct)

		req := httptest.NewRequest(http.MethodPost, "/api/products", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(CreateProductHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusConflict {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
		}

		response := decodeResponse(t, rr.Body)
		if response.Success {
			t.Error("Expected success to be false")
		}
	})
}

func TestCreateCategoryHandler(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		newCategory := models.Category{
			ID:          10,
			Name:        "Category 10",
			Description: "Desc 10",
		}
		body, _ := json.Marshal(newCategory)

		req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(CreateCategoryHandler)
		handler.ServeHTTP(rr, req)

		if status := rr.Code; status != http.StatusCreated {
			t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
		}

		response := decodeResponse(t, rr.Body)
		if !response.Success {
			t.Error("Expected success to be true")
		}
	})
}
