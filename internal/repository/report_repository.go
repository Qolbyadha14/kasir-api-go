package repository

import (
	"database/sql"
	"kasir-api-go/internal/models"
	"time"
)

type ReportRepository interface {
	GetSalesReport(startDate, endDate time.Time) (models.SalesReport, error)
}

type postgresReportRepository struct {
	db *sql.DB
}

func NewPostgresReportRepository(db *sql.DB) ReportRepository {
	return &postgresReportRepository{db: db}
}

func (r *postgresReportRepository) GetSalesReport(startDate, endDate time.Time) (models.SalesReport, error) {
	var report models.SalesReport

	// 1. Get total revenue and transactions
	query := `
		SELECT COALESCE(SUM(total_amount), 0), COUNT(id)
		FROM transactions
		WHERE created_at >= $1 AND created_at < $2`

	err := r.db.QueryRow(query, startDate, endDate).Scan(&report.TotalRevenue, &report.TotalTransactions)
	if err != nil {
		return report, err
	}

	// 2. Get best selling product
	bestSellingQuery := `
		SELECT p.name, COALESCE(SUM(td.quantity), 0) as total_qty
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		JOIN transactions t ON td.transaction_id = t.id
		WHERE t.created_at >= $1 AND t.created_at < $2
		GROUP BY p.name
		ORDER BY total_qty DESC
		LIMIT 1`

	err = r.db.QueryRow(bestSellingQuery, startDate, endDate).Scan(&report.BestSellingProduct.Name, &report.BestSellingProduct.QtySold)
	if err != nil && err != sql.ErrNoRows {
		return report, err
	}

	return report, nil
}
