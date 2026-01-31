package repository

import (
	"database/sql"
	"kasir-api-go/internal/models"
)

type PostgresProductRepository struct {
	db *sql.DB
}

func NewPostgresProductRepository(db *sql.DB) *PostgresProductRepository {
	return &PostgresProductRepository{db: db}
}

func (r *PostgresProductRepository) GetAll() []models.Product {
	query := `
		SELECT p.id, p.name, p.price, p.stock, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		ORDER BY p.id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Product{}
	}
	defer rows.Close()

	var products []models.Product
	for rows.Next() {
		var p models.Product
		var categoryID sql.NullInt64
		var categoryName, categoryDesc sql.NullString

		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName, &categoryDesc); err != nil {
			continue
		}

		if categoryID.Valid {
			p.Category = &models.Category{
				ID:          int(categoryID.Int64),
				Name:        categoryName.String,
				Description: categoryDesc.String,
			}
		}

		products = append(products, p)
	}

	return products
}

func (r *PostgresProductRepository) GetByID(id int) (models.Product, bool) {
	query := `
		SELECT p.id, p.name, p.price, p.stock, c.id, c.name, c.description
		FROM products p
		LEFT JOIN categories c ON p.category_id = c.id
		WHERE p.id = $1
	`

	var p models.Product
	var categoryID sql.NullInt64
	var categoryName, categoryDesc sql.NullString

	err := r.db.QueryRow(query, id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock, &categoryID, &categoryName, &categoryDesc)
	if err != nil {
		return models.Product{}, false
	}

	if categoryID.Valid {
		p.Category = &models.Category{
			ID:          int(categoryID.Int64),
			Name:        categoryName.String,
			Description: categoryDesc.String,
		}
	}

	return p, true
}

func (r *PostgresProductRepository) Create(product models.Product) {
	var categoryID *int
	if product.Category != nil {
		categoryID = &product.Category.ID
	}

	query := `INSERT INTO products (id, name, price, stock, category_id) VALUES ($1, $2, $3, $4, $5)`
	r.db.Exec(query, product.ID, product.Name, product.Price, product.Stock, categoryID)
}

func (r *PostgresProductRepository) Update(id int, product models.Product) bool {
	var categoryID *int
	if product.Category != nil {
		categoryID = &product.Category.ID
	}

	query := `UPDATE products SET name = $1, price = $2, stock = $3, category_id = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5`

	result, err := r.db.Exec(query, product.Name, product.Price, product.Stock, categoryID, id)
	if err != nil {
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

func (r *PostgresProductRepository) Delete(id int) bool {
	query := `DELETE FROM products WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}
