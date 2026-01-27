package repository

import "kasir-api-go/internal/models"

type ProductRepository interface {
	GetAll() []models.Product
	GetByID(id int) (models.Product, bool)
	Create(product models.Product)
	Update(id int, product models.Product) bool
	Delete(id int) bool
}

type InMemoryProductRepository struct {
	products []models.Product
}

func NewInMemoryProductRepository() *InMemoryProductRepository {
	// Re-creating the initial state
	c1 := models.Category{ID: 1, Name: "Category 1", Description: "Description 1"}
	c2 := models.Category{ID: 2, Name: "Category 2", Description: "Description 2"}

	return &InMemoryProductRepository{
		products: []models.Product{
			{
				ID:       1,
				Name:     "Product 1",
				Price:    10000,
				Stock:    10,
				Category: &c1,
			},
			{
				ID:       2,
				Name:     "Product 2",
				Price:    20000,
				Stock:    20,
				Category: &c2,
			},
		},
	}
}

func (r *InMemoryProductRepository) GetAll() []models.Product {
	return r.products
}

func (r *InMemoryProductRepository) GetByID(id int) (models.Product, bool) {
	for _, p := range r.products {
		if p.ID == id {
			return p, true
		}
	}
	return models.Product{}, false
}

func (r *InMemoryProductRepository) Create(product models.Product) {
	r.products = append(r.products, product)
}

func (r *InMemoryProductRepository) Update(id int, product models.Product) bool {
	for i, p := range r.products {
		if p.ID == id {
			r.products[i] = product
			return true
		}
	}
	return false
}

func (r *InMemoryProductRepository) Delete(id int) bool {
	for i, p := range r.products {
		if p.ID == id {
			r.products = append(r.products[:i], r.products[i+1:]...)
			return true
		}
	}
	return false
}
