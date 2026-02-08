package models

type Product struct {
	ID       int       `json:"id"`
	Name     string    `json:"name"`
	Price    int       `json:"price"`
	Stock    int       `json:"stock"`
	Category *Category `json:"category"`
}

type BestSellingProduct struct {
	Name    string `json:"nama"`
	QtySold int    `json:"qty_terjual"`
}

type SalesReport struct {
	TotalRevenue       int                `json:"total_revenue"`
	TotalTransactions  int                `json:"total_transaksi"`
	BestSellingProduct BestSellingProduct `json:"produk_terlaris"`
}
