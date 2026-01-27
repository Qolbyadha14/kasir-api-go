package handler

import (
	"encoding/json"
	"kasir-api-go/internal/api"
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/service"
	"net/http"
	"strconv"
	"strings"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{
		service: service,
	}
}

// @Summary List all categories
// @Description Get a list of all categories
// @Tags categories
// @Produce json
// @Success 200 {object} api.JSONResponse{data=[]models.Category}
// @Router /api/categories [get]
func (h *CategoryHandler) GetCategories(w http.ResponseWriter, r *http.Request) {
	categories := h.service.GetAll()
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
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	createdCategory, err := h.service.Create(category)
	if err != nil && err.Error() == "category ID already exists" {
		api.ErrorResponse(w, http.StatusConflict, err.Error(), "Duplicate ID")
		return
	}

	if err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, err.Error(), "Bad Request")
		return
	}

	api.SuccessResponse(w, http.StatusCreated, "Category created successfully", createdCategory)
}

// @Summary Get a category detail
// @Description Get details of a category by ID
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} api.JSONResponse{data=models.Category}
// @Failure 404 {object} api.JSONResponse
// @Router /api/categories/{id} [get]
func (h *CategoryHandler) GetCategoryDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	category, err := h.service.GetByID(id)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
		return
	}
	api.SuccessResponse(w, http.StatusOK, "Success", category)
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
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	var category models.Category
	if err := json.NewDecoder(r.Body).Decode(&category); err != nil {
		api.ErrorResponse(w, http.StatusBadRequest, "Invalid request body", err.Error())
		return
	}

	category.ID = id
	updatedCategory, err := h.service.Update(id, category)
	if err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
		return
	}
	api.SuccessResponse(w, http.StatusOK, "Category updated successfully", updatedCategory)
}

// @Summary Delete a category
// @Description Remove a category from the catalog
// @Tags categories
// @Produce json
// @Param id path int true "Category ID"
// @Success 200 {object} api.JSONResponse
// @Failure 404 {object} api.JSONResponse
// @Router /api/categories/{id} [delete]
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/categories/")
	id, _ := strconv.Atoi(idStr)

	if err := h.service.Delete(id); err != nil {
		api.ErrorResponse(w, http.StatusNotFound, "Category not found", "Category not found")
		return
	}
	api.SuccessResponse(w, http.StatusOK, "Category deleted successfully", nil)
}
