package repository

import (
	"database/sql"
	"kasir-api-go/internal/models"
)

type PostgresCategoryRepository struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) *PostgresCategoryRepository {
	return &PostgresCategoryRepository{db: db}
}

func (r *PostgresCategoryRepository) GetAll() []models.Category {
	query := `SELECT id, name, description FROM categories ORDER BY id`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.Category{}
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var c models.Category
		if err := rows.Scan(&c.ID, &c.Name, &c.Description); err != nil {
			continue
		}
		categories = append(categories, c)
	}

	return categories
}

func (r *PostgresCategoryRepository) GetByID(id int) (models.Category, bool) {
	query := `SELECT id, name, description FROM categories WHERE id = $1`

	var c models.Category
	err := r.db.QueryRow(query, id).Scan(&c.ID, &c.Name, &c.Description)
	if err != nil {
		return models.Category{}, false
	}

	return c, true
}

func (r *PostgresCategoryRepository) Create(category models.Category) {
	query := `INSERT INTO categories (id, name, description) VALUES ($1, $2, $3)`
	r.db.Exec(query, category.ID, category.Name, category.Description)
}

func (r *PostgresCategoryRepository) Update(id int, category models.Category) bool {
	query := `UPDATE categories SET name = $1, description = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`

	result, err := r.db.Exec(query, category.Name, category.Description, id)
	if err != nil {
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}

func (r *PostgresCategoryRepository) Delete(id int) bool {
	query := `DELETE FROM categories WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return false
	}

	rowsAffected, _ := result.RowsAffected()
	return rowsAffected > 0
}
