package service

import (
	"errors"
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/repository"
)

type ProductService interface {
	GetAll() []models.Product
	GetByID(id int) (models.Product, error)
	Create(product models.Product) (models.Product, error)
	Update(id int, product models.Product) (models.Product, error)
	Delete(id int) error
}

type productService struct {
	productRepo  repository.ProductRepository
	categoryRepo repository.CategoryRepository
}

func NewProductService(productRepo repository.ProductRepository, categoryRepo repository.CategoryRepository) ProductService {
	return &productService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productService) GetAll() []models.Product {
	return s.productRepo.GetAll()
}

func (s *productService) GetByID(id int) (models.Product, error) {
	product, found := s.productRepo.GetByID(id)
	if !found {
		return models.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) Create(product models.Product) (models.Product, error) {
	// Validation: Duplicate ID
	if _, found := s.productRepo.GetByID(product.ID); found {
		return models.Product{}, errors.New("product ID already exists")
	}

	// Validation: Category existence
	if product.Category != nil {
		if _, found := s.categoryRepo.GetByID(product.Category.ID); !found {
			return models.Product{}, errors.New("category not found")
		}
	}

	s.productRepo.Create(product)
	return product, nil
}

func (s *productService) Update(id int, product models.Product) (models.Product, error) {
	// Validation: Category existence
	if product.Category != nil {
		if _, found := s.categoryRepo.GetByID(product.Category.ID); !found {
			return models.Product{}, errors.New("category not found")
		}
	}

	if ok := s.productRepo.Update(id, product); !ok {
		return models.Product{}, errors.New("product not found")
	}
	return product, nil
}

func (s *productService) Delete(id int) error {
	if ok := s.productRepo.Delete(id); !ok {
		return errors.New("product not found")
	}
	return nil
}
