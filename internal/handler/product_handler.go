package handler

import (
	"encoding/json"
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/service"
	"kasir-api-go/internal/utils"
	"net/http"
	"strconv"
	"strings"
)

type ProductHandler struct {
	service service.ProductService
}

func NewProductHandler(service service.ProductService) *ProductHandler {
	return &ProductHandler{
		service: service,
	}
}

// @Summary List all products
// @Description Get a list of all products, optionally filtered by name
// @Tags products
// @Produce json
// @Param search query string false "Search products by name"
// @Success 200 {object} utils.JSONResponse{data=[]models.Product}
// @Router /api/products [get]
func (h *ProductHandler) GetProducts(w http.ResponseWriter, r *http.Request) {
	search := r.URL.Query().Get("search")
	products := h.service.GetAll(search)
	utils.SuccessResponse(w, http.StatusOK, "Success", products)
}

// @Summary Create a new product
// @Description Add a new product to the catalog
// @Tags products
// @Accept json
// @Produce json
// @Param product body models.Product true "Product object"
// @Success 201 {object} utils.JSONResponse{data=models.Product}
// @Failure 400 {object} utils.JSONResponse
// @Failure 409 {object} utils.JSONResponse
// @Router /api/products [post]
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	createdProduct, err := h.service.Create(product)
	if err != nil && err.Error() == "product ID already exists" {
		utils.ErrorResponse(w, http.StatusConflict, err.Error(), "Duplicate ID")
		return
	}

	if err != nil && err.Error() == "category not found" {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid Category ID")
		return
	}

	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Bad Request")
		return
	}

	utils.SuccessResponse(w, http.StatusCreated, "Product created successfully", createdProduct)
}

// @Summary Get a product detail
// @Description Get details of a product by ID
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.JSONResponse{data=models.Product}
// @Failure 404 {object} utils.JSONResponse
// @Router /api/products/{id} [get]
func (h *ProductHandler) GetProductDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	product, err := h.service.GetByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Success", product)
}

// @Summary Update a product
// @Description Update an existing product's details
// @Tags products
// @Accept json
// @Produce json
// @Param id path int true "Product ID"
// @Param product body models.Product true "Product object"
// @Success 200 {object} utils.JSONResponse{data=models.Product}
// @Failure 400 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Router /api/products/{id} [put]
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	var product models.Product
	if err := json.NewDecoder(r.Body).Decode(&product); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	product.ID = id
	updatedProduct, err := h.service.Update(id, product)
	if err != nil && err.Error() == "category not found" {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Invalid Category ID")
		return
	}

	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Product updated successfully", updatedProduct)
}

// @Summary Delete a product
// @Description Remove a product from the catalog
// @Tags products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} utils.JSONResponse
// @Failure 404 {object} utils.JSONResponse
// @Router /api/products/{id} [delete]
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/products/")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.Delete(id); err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Product not found", "Product not found")
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Product deleted successfully", nil)
}
