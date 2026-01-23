package main

import (
	"encoding/json"
	"fmt"
	"kasir-api-go/internal/api"
	"kasir-api-go/internal/models"
	"net/http"
	"strconv"
	"strings"

	_ "kasir-api-go/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Use models from internal package
type Category = models.Category
type Product = models.Product

var categories = []models.Category{
	{
		ID:          1,
		Name:        "Category 1",
		Description: "Description 1",
	},
	{
		ID:          2,
		Name:        "Category 2",
		Description: "Description 2",
	},
}

var products = []models.Product{
	{
		ID:       1,
		Name:     "Product 1",
		Price:    10000,
		Stock:    10,
		Category: &categories[0],
	},
	{
		ID:       2,
		Name:     "Product 2",
		Price:    20000,
		Stock:    20,
		Category: &categories[1],
	},
}

func getProduct(id int) (models.Product, bool) {
	for _, p := range products {
		if p.ID == id {
			return p, true
		}
	}
	return models.Product{}, false
}

func updateProduct(id int, product models.Product) bool {
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

func getCategory(id int) (models.Category, bool) {
	for _, c := range categories {
		if c.ID == id {
			return c, true
		}
	}
	return models.Category{}, false
}

func updateCategory(id int, category models.Category) bool {
	for i, c := range categories {
		if c.ID == id {
			categories[i] = category
			return true
		}
	}
	return false
}

func deleteCategory(id int) bool {
	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			return true
		}
	}
	return false
}

// @Summary Health check
// @Description Get the status of the API
// @Tags health
// @Produce json
// @Success 200 {object} api.JSONResponse
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	api.SuccessResponse(w, http.StatusOK, "API Running", map[string]string{"status": "ok"})
}

// @Summary List all products
// @Description Get a list of all products
// @Tags products
// @Produce json
// @Success 200 {object} api.JSONResponse{data=[]models.Product}
// @Router /api/products [get]
func GetProductsHandler(w http.ResponseWriter, r *http.Request) {
	api.SuccessResponse(w, http.StatusOK, "Success", products)
}

// @Summary Create a new product
// @Description Add a new product to the catalog
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} api.JSONResponse{data=models.Product}
// @Failure 400 {object} api.JSONResponse
// @Failure 409 {object} api.JSONResponse
// @Router /api/products [post]
func CreateProductHandler(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validation: Duplicate ID
	if _, found := getProduct(product.ID); found {
		api.ErrorResponse(w, http.StatusConflict, "Product ID already exists", "Duplicate ID")
		return
	}

	// Validation: Category existence
	if product.Category != nil {
		if _, found := getCategory(product.Category.ID); !found {
			api.ErrorResponse(w, http.StatusBadRequest, "Category not found", "Invalid Category ID")
			return
		}
	}

	products = append(products, product)
	api.SuccessResponse(w, http.StatusCreated, "Product created successfully", product)
}

// @Summary List all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {object} api.JSONResponse{data=[]models.Category}
// @Router /api/categories [get]
func GetCategoriesHandler(w http.ResponseWriter, r *http.Request) {
	api.SuccessResponse(w, http.StatusOK, "Success", categories)
}

// @Summary Create a new category
// @Description Add a new category to the catalog
// @Tags categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category object"
// @Success 201 {object} api.JSONResponse{data=models.Category}
// @Failure 400 {object} api.JSONResponse
// @Failure 409 {object} api.JSONResponse
// @Router /api/categories [post]
func CreateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validation: Duplicate ID
	if _, found := getCategory(category.ID); found {
		api.ErrorResponse(w, http.StatusConflict, "Category ID already exists", "Duplicate ID")
		return
	}

	categories = append(categories, category)
	api.SuccessResponse(w, http.StatusCreated, "Category created successfully", category)
}

// @Summary Get a category detail
// @Description Get details of a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} api.JSONResponse{data=models.Category}
// @Failure 404 {object} api.JSONResponse
// @Router /api/categories/{id} [get]
func GetCategoryDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	if category, found := getCategory(id); found {
		api.SuccessResponse(w, http.StatusOK, "Success", category)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
}

// @Summary Update a category
// @Description Update an existing category's details
// @Tags categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Category object"
// @Success 200 {object} api.JSONResponse{data=models.Category}
// @Failure 400 {object} api.JSONResponse
// @Failure 404 {object} api.JSONResponse
// @Router /api/categories/{id} [put]
func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category.ID = id
	if ok := updateCategory(id, category); ok {
		api.SuccessResponse(w, http.StatusOK, "Category updated successfully", category)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
}

// @Summary Delete a category
// @Description Remove a category from the catalog
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} api.JSONResponse
// @Failure 404 {object} api.JSONResponse
// @Router /api/categories/{id} [delete]
func DeleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	if ok := deleteCategory(id); ok {
		api.SuccessResponse(w, http.StatusOK, "Category deleted successfully", nil)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
}

// @Summary Get a product detail
// @Description Get details of a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} api.JSONResponse{data=models.Product}
// @Failure 404 {object} api.JSONResponse
// @Router /api/products/{id} [get]
func GetProductDetailHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	if product, found := getProduct(id); found {
		api.SuccessResponse(w, http.StatusOK, "Success", product)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
}

// @Summary Update a product
// @Description Update an existing product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} api.JSONResponse{data=models.Product}
// @Failure 400 {object} api.JSONResponse
// @Failure 404 {object} api.JSONResponse
// @Router /api/products/{id} [put]
func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	// Validation: Category existence
	if product.Category != nil {
		if _, found := getCategory(product.Category.ID); !found {
			api.ErrorResponse(w, http.StatusBadRequest, "Category not found", "Invalid Category ID")
			return
		}
	}

	product.ID = id
	if ok := updateProduct(id, product); ok {
		api.SuccessResponse(w, http.StatusOK, "Product updated successfully", product)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
}

// @Summary Delete a product
// @Description Remove a product from the catalog
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} api.JSONResponse
// @Failure 404 {object} api.JSONResponse
// @Router /api/products/{id} [delete]
func DeleteProductHandler(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	if ok := deleteProduct(id); ok {
		api.SuccessResponse(w, http.StatusOK, "Product deleted successfully", nil)
		return
	}
	api.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
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
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateProductHandler(w, r)
		} else {
			GetProductsHandler(w, r)
		}
	})

	// Handle /api/products/{id} (GET, UPDATE AND DELETE)
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
		if idStr == "" {
			return
		}
		switch r.Method {
		case http.MethodGet:
			GetProductDetailHandler(w, r)
		case http.MethodPut:
			UpdateProductHandler(w, r)
		case http.MethodDelete:
			DeleteProductHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Handle /api/categories (GET and POST)
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			CreateCategoryHandler(w, r)
		} else {
			GetCategoriesHandler(w, r)
		}
	})

	// Handle /api/categories/{id} (GET, UPDATE AND DELETE)
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		if idStr == "" {
			return
		}
		switch r.Method {
		case http.MethodGet:
			GetCategoryDetailHandler(w, r)
		case http.MethodPut:
			UpdateCategoryHandler(w, r)
		case http.MethodDelete:
			DeleteCategoryHandler(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Serve static files from the "public" directory
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Swagger UI
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Println("Starting server on http://localhost:8080")
	fmt.Println("Swagger documentation at http://localhost:8080/swagger/index.html")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
