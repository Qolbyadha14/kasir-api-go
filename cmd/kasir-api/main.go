package main

import (
	"fmt"
	"net/http"
	"strings"

	"kasir-api-go/docs"
	"kasir-api-go/internal/config"
	"kasir-api-go/internal/database"
	"kasir-api-go/internal/handler"
	"kasir-api-go/internal/repository"
	"kasir-api-go/internal/service"

	_ "kasir-api-go/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

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
// @schemes http https
func main() {
	// Load configuration
	cfg := config.GetConfig()

	// Database initialization
	db, closeDB, err := database.NewPostgres(&cfg.Database)
	if err != nil {
		fmt.Println("Database connection failed:", err)
		return
	}
	defer closeDB()
	fmt.Println("Database connected successfully")

	port := cfg.App.Port

	// Dependency Injection
	// Repositories
	categoryRepo := repository.NewPostgresCategoryRepository(db)
	productRepo := repository.NewPostgresProductRepository(db)

	// Update swagger info host dynamically
	if cfg.App.URL != "" {
		docs.SwaggerInfo.Host = cfg.App.URL
	}

	// Services
	categoryService := service.NewCategoryService(categoryRepo)
	productService := service.NewProductService(productRepo, categoryRepo)

	// Handlers
	categoryHandler := handler.NewCategoryHandler(categoryService)
	productHandler := handler.NewProductHandler(productService)

	// Get localhost:8080/health
	http.HandleFunc("/health", handler.HealthHandler)

	// Handle /api/products (GET and POST)
	http.HandleFunc("/api/products", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			productHandler.CreateProduct(w, r)
			return
		}
		productHandler.GetProducts(w, r)
	})

	// Handle /api/products/{id} (GET, UPDATE AND DELETE)
	http.HandleFunc("/api/products/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
		if idStr == "" {
			return
		}
		switch r.Method {
		case http.MethodGet:
			productHandler.GetProductDetail(w, r)
		case http.MethodPut:
			productHandler.UpdateProduct(w, r)
		case http.MethodDelete:
			productHandler.DeleteProduct(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Handle /api/categories (GET and POST)
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			categoryHandler.CreateCategory(w, r)
			return
		}
		categoryHandler.GetCategories(w, r)
	})

	// Handle /api/categories/{id} (GET, UPDATE AND DELETE)
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
		if idStr == "" {
			return
		}
		switch r.Method {
		case http.MethodGet:
			categoryHandler.GetCategoryDetail(w, r)
		case http.MethodPut:
			categoryHandler.UpdateCategory(w, r)
		case http.MethodDelete:
			categoryHandler.DeleteCategory(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	// Serve static files from the "public" directory
	// Assuming the app is run from the project root
	fs := http.FileServer(http.Dir("./public"))
	http.Handle("/", fs)

	// Swagger UI
	http.Handle("/swagger/", httpSwagger.WrapHandler)

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	fmt.Printf("Swagger documentation at http://localhost:%s/swagger/index.html\n", port)

	err = http.ListenAndServe(":"+port, nil)

	if err != nil {
		fmt.Println("error starting server:", err)
	}
}
