package repository

import "kasir-api-go/internal/models"

type CategoryRepository interface {
	GetAll() []models.Category
	GetByID(id int) (models.Category, bool)
	Create(category models.Category)
	Update(id int, category models.Category) bool
	Delete(id int) bool
}

type InMemoryCategoryRepository struct {
	categories []models.Category
}

func NewInMemoryCategoryRepository() *InMemoryCategoryRepository {
	return &InMemoryCategoryRepository{
		categories: []models.Category{
			{
				ID:          1,
				Name:        "Category 1",
				Description: "Description 1",
			},
			{
				ID:          2,
				Name:        "Category 2",
				Description: "Description 2",
			},
		},
	}
}

func (r *InMemoryCategoryRepository) GetAll() []models.Category {
	return r.categories
}

func (r *InMemoryCategoryRepository) GetByID(id int) (models.Category, bool) {
	for _, c := range r.categories {
		if c.ID == id {
			return c, true
		}
	}
	return models.Category{}, false
}

func (r *InMemoryCategoryRepository) Create(category models.Category) {
	r.categories = append(r.categories, category)
}

func (r *InMemoryCategoryRepository) Update(id int, category models.Category) bool {
	for i, c := range r.categories {
		if c.ID == id {
			r.categories[i] = category
			return true
		}
	}
	return false
}

func (r *InMemoryCategoryRepository) Delete(id int) bool {
	for i, c := range r.categories {
		if c.ID == id {
			r.categories = append(r.categories[:i], r.categories[i+1:]...)
			return true
		}
	}
	return false
}
