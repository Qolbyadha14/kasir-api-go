package service

import (
	"errors"
	"kasir-api-go/internal/models"
	"kasir-api-go/internal/repository"
)

type TransactionService interface {
	Checkout(items []models.CheckoutItem) (models.Transaction, error)
	GetAllTransactions() ([]models.Transaction, error)
	GetTransactionByID(id int) (models.Transaction, error)
}

type transactionService struct {
	repo        repository.TransactionRepository
	productRepo repository.ProductRepository
}

func NewTransactionService(repo repository.TransactionRepository, productRepo repository.ProductRepository) TransactionService {
	return &transactionService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *transactionService) Checkout(items []models.CheckoutItem) (models.Transaction, error) {
	if len(items) == 0 {
		return models.Transaction{}, errors.New("transaction items cannot be empty")
	}

	transaction, err := s.repo.CreateTransaction(items)
	if err != nil {
		return models.Transaction{}, err
	}

	return *transaction, nil
}

func (s *transactionService) GetAllTransactions() ([]models.Transaction, error) {
	return s.repo.GetAll()
}

func (s *transactionService) GetTransactionByID(id int) (models.Transaction, error) {
	return s.repo.GetByID(id)
}
