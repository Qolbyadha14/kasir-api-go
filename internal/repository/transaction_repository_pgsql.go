package repository

import (
	"database/sql"
	"fmt"
	"kasir-api-go/internal/models"
	"sort"
	"strings"
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
	if len(items) == 0 {
		return nil, fmt.Errorf("transaction items cannot be empty")
	}

	// 1. Consolidate duplicate products
	consolidated := make(map[int]int)
	for _, item := range items {
		consolidated[item.ProductID] += item.Quantity
	}

	// 2. Sort Product IDs to prevent deadlocks
	productIDs := make([]int, 0, len(consolidated))
	for id := range consolidated {
		productIDs = append(productIDs, id)
	}
	sort.Ints(productIDs)

	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	// 3. Process products in sorted order with locking
	for _, id := range productIDs {
		qty := consolidated[id]
		var productPrice, stock int
		var productName string

		// Use FOR UPDATE to lock the row and prevent race conditions
		err := tx.QueryRow("SELECT name, price, stock FROM products WHERE id = $1 FOR UPDATE", id).
			Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("product id %d not found", id)
		}
		if err != nil {
			return nil, err
		}

		if stock < qty {
			return nil, fmt.Errorf("insufficient stock for product: %s", productName)
		}

		subtotal := productPrice * qty
		totalAmount += subtotal

		_, err = tx.Exec("UPDATE products SET stock = stock - $1 WHERE id = $2", qty, id)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   id,
			ProductName: productName,
			Quantity:    qty,
			Subtotal:    subtotal,
		})
	}

	// 4. Insert transaction header
	var transactionID int
	err = tx.QueryRow("INSERT INTO transactions (total_amount) VALUES ($1) RETURNING id", totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	// 5. Bulk insert transaction details
	if len(details) > 0 {
		valueStrings := make([]string, 0, len(details))
		valueArgs := make([]interface{}, 0, len(details)*4)
		for i, d := range details {
			pos := i * 4
			valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", pos+1, pos+2, pos+3, pos+4))
			valueArgs = append(valueArgs, transactionID, d.ProductID, d.Quantity, d.Subtotal)
		}
		bulkInsertQuery := fmt.Sprintf("INSERT INTO transaction_details (transaction_id, product_id, quantity, subtotal) VALUES %s",
			strings.Join(valueStrings, ","))

		_, err = tx.Exec(bulkInsertQuery, valueArgs...)
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
	var ids []interface{}
	transMap := make(map[int]*models.Transaction)

	for rows.Next() {
		var t models.Transaction
		if err := rows.Scan(&t.ID, &t.TotalAmount, &t.CreatedAt); err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
		ids = append(ids, t.ID)
		transMap[t.ID] = &transactions[len(transactions)-1]
	}

	if len(ids) == 0 {
		return transactions, nil
	}

	// Fetch all details for these transactions in one go
	detailQuery := `
		SELECT td.id, td.transaction_id, td.product_id, td.quantity, td.subtotal, p.name, p.price
		FROM transaction_details td
		JOIN products p ON td.product_id = p.id
		WHERE td.transaction_id IN (`

	for i := range ids {
		detailQuery += fmt.Sprintf("$%d", i+1)
		if i < len(ids)-1 {
			detailQuery += ","
		}
	}
	detailQuery += ")"

	detailRows, err := r.db.Query(detailQuery, ids...)
	if err != nil {
		return nil, err
	}
	defer detailRows.Close()

	for detailRows.Next() {
		var d models.TransactionDetail
		var p models.Product
		if err := detailRows.Scan(&d.ID, &d.TransactionID, &d.ProductID, &d.Quantity, &d.Subtotal, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		p.ID = d.ProductID
		d.Product = &p
		d.ProductName = p.Name

		if t, ok := transMap[d.TransactionID]; ok {
			t.Details = append(t.Details, d)
		}
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
