package models

import "time"

type Transaction struct {
	ID          int                 `json:"id"`
	TotalAmount int                 `json:"total_amount"`
	CreatedAt   time.Time           `json:"created_at"`
	Details     []TransactionDetail `json:"details,omitempty"`
}

type TransactionDetail struct {
	ID            int      `json:"id"`
	TransactionID int      `json:"transaction_id"`
	ProductID     int      `json:"product_id"`
	Product       *Product `json:"product,omitempty"`
	Quantity      int      `json:"quantity"`
	Subtotal      int      `json:"subtotal"`
}
