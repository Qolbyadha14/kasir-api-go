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

type TransactionHandler struct {
	service service.TransactionService
}

func NewTransactionHandler(service service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

// @Summary List all transactions
// @Description Get a list of all transactions including their details
// @Tags transactions
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=[]models.Transaction}
// @Router /api/transactions [get]
func (h *TransactionHandler) GetTransactions(w http.ResponseWriter, r *http.Request) {
	transactions, err := h.service.GetAllTransactions()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch transactions", err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Success", transactions)
}

// @Summary Checkout transactions
// @Description Create a new transaction from multiple items and update stock
// @Tags transactions
// @Accept json
// @Produce json
// @Param request body models.CheckoutRequest true "Checkout Request object"
// @Success 201 {object} models.Transaction
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /api/checkout [post]
func (h *TransactionHandler) HandleCheckout(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		h.Checkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *TransactionHandler) Checkout(w http.ResponseWriter, r *http.Request) {
	var req models.CheckoutRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	transaction, err := h.service.Checkout(req.Items)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)
}

// @Summary Get a transaction detail
// @Description Get details of a transaction by ID
// @Tags transactions
// @Produce json
// @Param id path int true "Transaction ID"
// @Success 200 {object} utils.JSONResponse{data=models.Transaction}
// @Failure 404 {object} utils.JSONResponse
// @Router /api/transactions/{id} [get]
func (h *TransactionHandler) GetTransactionDetail(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/transactions/")
	id, _ := strconv.Atoi(idStr)

	transaction, err := h.service.GetTransactionByID(id)
	if err != nil {
		utils.ErrorResponse(w, http.StatusNotFound, "Transaction not found", "Transaction not found")
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Success", transaction)
}
