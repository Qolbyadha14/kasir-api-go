package handler

import (
	"kasir-api-go/internal/utils"
	"net/http"
)

// @Summary Health check
// @Description Get the status of the API
// @Tags health
// @Produce json
// @Success 200 {object} utils.JSONResponse
// @Router /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	utils.SuccessResponse(w, http.StatusOK, "API Running", map[string]string{"status": "ok"})
}
