package service

import (
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/repository"
	"time"
)

type ReportService interface {
	GetTodayReport() (models.SalesReport, error)
	GetReportByRange(startDate, endDate time.Time) (models.SalesReport, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo: repo}
}

func (s *reportService) GetTodayReport() (models.SalesReport, error) {
	now := time.Now()
	startDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endDate := startDate.Add(24 * time.Hour)
	return s.repo.GetSalesReport(startDate, endDate)
}

func (s *reportService) GetReportByRange(startDate, endDate time.Time) (models.SalesReport, error) {
	// Adjust endDate to include the whole day if only date is provided
	endDate = endDate.Add(24 * time.Hour)
	return s.repo.GetSalesReport(startDate, endDate)
}
