package handler

import (
	"kasir-api-go/internal/service"
	"kasir-api-go/internal/utils"
	"net/http"
	"time"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(service service.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

// @Summary Get sales report for today
// @Description Get total revenue, total transactions, and best selling product for today
// @Tags report
// @Produce json
// @Success 200 {object} utils.JSONResponse{data=models.SalesReport}
// @Router /api/report/today [get]
func (h *ReportHandler) GetTodayReport(w http.ResponseWriter, r *http.Request) {
	report, err := h.service.GetTodayReport()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch today's report", err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Success", report)
}

// @Summary Get sales report by date range
// @Description Get total revenue, total transactions, and best selling product for a specific date range
// @Tags report
// @Produce json
// @Param start_date query string true "Start date (YYYY-MM-DD)"
// @Param end_date query string true "End date (YYYY-MM-DD)"
// @Success 200 {object} utils.JSONResponse{data=models.SalesReport}
// @Router /api/report [get]
func (h *ReportHandler) GetReportByRange(w http.ResponseWriter, r *http.Request) {
	startDateStr := r.URL.Query().Get("start_date")
	endDateStr := r.URL.Query().Get("end_date")

	if startDateStr == "" || endDateStr == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid date range", "start_date and end_date are required")
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid start_date format", "Expected YYYY-MM-DD")
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid end_date format", "Expected YYYY-MM-DD")
		return
	}

	report, err := h.service.GetReportByRange(startDate, endDate)
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, "Failed to fetch report", err.Error())
		return
	}
	utils.SuccessResponse(w, http.StatusOK, "Success", report)
}
