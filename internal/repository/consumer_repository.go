package repository

import (
	"main/internal/model"

	"gorm.io/gorm"
)

// ConsumerRepository defines all operations for Consumer entity
type ConsumerRepository interface {
	Create(consumer *model.Consumer) error
	GetByID(id uint) (*model.Consumer, error)
	GetByNIK(nik string) (*model.Consumer, error)
	GetAll() ([]model.Consumer, error)
	Update(consumer *model.Consumer) error
	Delete(id uint) error
}

// ConsumerLimitRepository defines all operations for ConsumerLimit entity
type ConsumerLimitRepository interface {
	Create(limit *model.ConsumerLimit) error
	GetByID(id uint) (*model.ConsumerLimit, error)
	GetByConsumerAndTenor(consumerID uint, tenor int) (*model.ConsumerLimit, error)
	GetByConsumerID(consumerID uint) ([]model.ConsumerLimit, error)
	Update(limit *model.ConsumerLimit) error
	Delete(id uint) error
}

// TransactionRepository defines all operations for Transaction entity
type TransactionRepository interface {
	Create(transaction *model.Transaction) error
	GetByID(id uint) (*model.Transaction, error)
	GetByContractNumber(contractNumber string) (*model.Transaction, error)
	GetByConsumerID(consumerID uint) ([]model.Transaction, error)
	Update(transaction *model.Transaction) error
	Delete(id uint) error
}

// consumerRepository is the implementation of ConsumerRepository
type consumerRepository struct {
	db *gorm.DB
}

// NewConsumerRepository creates a new instance of ConsumerRepository
func NewConsumerRepository(db *gorm.DB) ConsumerRepository {
	return &consumerRepository{db: db}
}

func (r *consumerRepository) Create(consumer *model.Consumer) error {
	return r.db.Create(consumer).Error
}

func (r *consumerRepository) GetByID(id uint) (*model.Consumer, error) {
	var consumer model.Consumer
	err := r.db.Where("id = ?", id).First(&consumer).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (r *consumerRepository) GetByNIK(nik string) (*model.Consumer, error) {
	var consumer model.Consumer
	err := r.db.Where("nik = ?", nik).First(&consumer).Error
	if err != nil {
		return nil, err
	}
	return &consumer, nil
}

func (r *consumerRepository) GetAll() ([]model.Consumer, error) {
	var consumers []model.Consumer
	err := r.db.Find(&consumers).Error
	return consumers, err
}

func (r *consumerRepository) Update(consumer *model.Consumer) error {
	return r.db.Save(consumer).Error
}

func (r *consumerRepository) Delete(id uint) error {
	return r.db.Delete(&model.Consumer{}, id).Error
}

// consumerLimitRepository is the implementation of ConsumerLimitRepository
type consumerLimitRepository struct {
	db *gorm.DB
}

// NewConsumerLimitRepository creates a new instance of ConsumerLimitRepository
func NewConsumerLimitRepository(db *gorm.DB) ConsumerLimitRepository {
	return &consumerLimitRepository{db: db}
}

func (r *consumerLimitRepository) Create(limit *model.ConsumerLimit) error {
	return r.db.Create(limit).Error
}

func (r *consumerLimitRepository) GetByID(id uint) (*model.ConsumerLimit, error) {
	var limit model.ConsumerLimit
	err := r.db.Where("id = ?", id).First(&limit).Error
	return &limit, err
}

func (r *consumerLimitRepository) GetByConsumerAndTenor(consumerID uint, tenor int) (*model.ConsumerLimit, error) {
	var limit model.ConsumerLimit
	err := r.db.Where("consumer_id = ? AND tenor = ?", consumerID, tenor).First(&limit).Error
	return &limit, err
}

func (r *consumerLimitRepository) GetByConsumerID(consumerID uint) ([]model.ConsumerLimit, error) {
	var limits []model.ConsumerLimit
	err := r.db.Where("consumer_id = ?", consumerID).Find(&limits).Error
	return limits, err
}

func (r *consumerLimitRepository) Update(limit *model.ConsumerLimit) error {
	return r.db.Save(limit).Error
}

func (r *consumerLimitRepository) Delete(id uint) error {
	return r.db.Delete(&model.ConsumerLimit{}, id).Error
}

// transactionRepository is the implementation of TransactionRepository
type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository creates a new instance of TransactionRepository
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetByID(id uint) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("id = ?", id).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByContractNumber(contractNumber string) (*model.Transaction, error) {
	var transaction model.Transaction
	err := r.db.Where("contract_number = ?", contractNumber).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) GetByConsumerID(consumerID uint) ([]model.Transaction, error) {
	var transactions []model.Transaction
	err := r.db.Where("consumer_id = ?", consumerID).Find(&transactions).Error
	return transactions, err
}

func (r *transactionRepository) Update(transaction *model.Transaction) error {
	return r.db.Save(transaction).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Transaction{}, id).Error
}
