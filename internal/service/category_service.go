package service

import (
	"errors"
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/repository"
)

type CategoryService interface {
	GetAll() []models.Category
	GetByID(id int) (models.Category, error)
	Create(category models.Category) (models.Category, error)
	Update(id int, category models.Category) (models.Category, error)
	Delete(id int) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{
		repo: repo,
	}
}

func (s *categoryService) GetAll() []models.Category {
	return s.repo.GetAll()
}

func (s *categoryService) GetByID(id int) (models.Category, error) {
	category, found := s.repo.GetByID(id)
	if !found {
		return models.Category{}, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) Create(category models.Category) (models.Category, error) {
	// Validation: Duplicate ID check
	if _, found := s.repo.GetByID(category.ID); found {
		return models.Category{}, errors.New("category ID already exists")
	}

	s.repo.Create(category)
	return category, nil
}

func (s *categoryService) Update(id int, category models.Category) (models.Category, error) {
	if ok := s.repo.Update(id, category); !ok {
		return models.Category{}, errors.New("category not found")
	}
	return category, nil
}

func (s *categoryService) Delete(id int) error {
	if ok := s.repo.Delete(id); !ok {
		return errors.New("category not found")
	}
	return nil
}
