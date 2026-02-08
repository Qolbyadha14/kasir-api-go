package repository

import (
	"database/sql"
	"fmt"
	"kasir-api-go/internal/models"
)

type TransactionRepository interface {
	CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error)
	GetAll() ([]models.Transaction, error)
	GetByID(id int) (models.Transaction, error)
}

type postgresTransactionRepository struct {
	db *sql.DB
}

func NewPostgresTransactionRepository(db *sql.DB) TransactionRepository {
	return &postgresTransactionRepository{db: db}
}

func (r *postgresTransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1", item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		if stock < item.Quantity {
			return nil, fmt.Errorf("insufficient stock for product: %s", productName)
		}

		subtotal := productPrice * item.Quantity
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subtotal,
		})
	}

	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	for i := range details {
		details[i].TransactionID = transactionID
		_, err = tx.Exec("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES ($1, $2, $3, $4)",
			transactionID, details[i].ProductID, details[i].Quantity, details[i].Subtotal)
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}

func (r *postgresTransactionRepository) GetAll() ([]models.Transaction, error) {
	query := `SELECT id, total_amount, created_at FROM transactions ORDER BY created_at DESC`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}
	return transactions, nil
}

func (r *postgresTransactionRepository) GetByID(id int) (models.Transaction, error) {
	var t models.Transaction
	query := `SELECT id, total_amount, created_at FROM transactions WHERE id = $1`
	err := r.db.QueryRow(query, id).Scan(&t.ID, &t.TotalAmount, &t.CreatedAt)
	if err != nil {
		return t, err
	}

	detailQuery := `
		SELECT td.id, td.transaction_id, td.product_id, td.quantity, td.subtotal, p.name, p.price
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id = $1`

	rows, err := r.db.Query(detailQuery, t.ID)
	if err != nil {
		return t, err
	}
	defer rows.Close()

	for rows.Next() {
		var d models.TransactionDetail
		var p models.Product
		if err := rows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.Quantity, &d.Subtotal, &p.Name, &p.Price); err != nil {
			return t, err
		}
		p.ID = d.ProductID
		d.Product = &p
		d.ProductName = p.Name
		t.Details = append(t.Details, d)
	}

	return t, nil
}
