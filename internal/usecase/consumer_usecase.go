package usecase

import (
	"errors"
	"log"
	"regexp"
	"sync"
	"time"

	"main/internal/model"
	"main/internal/repository"
	"gorm.io/gorm"
)

// ConsumerUsecase defines all business logic operations for Consumer
type ConsumerUsecase interface {
	RegisterConsumer(consumer *model.Consumer) error
	GetConsumer(id uint) (*model.Consumer, error)
	GetConsumerByNIK(nik string) (*model.Consumer, error)
	UpdateConsumer(consumer *model.Consumer) error
	DeleteConsumer(id uint) error
}

// ConsumerLimitUsecase defines all business logic operations for ConsumerLimit
type ConsumerLimitUsecase interface {
	AssignLimit(limit *model.ConsumerLimit) error
	GetLimitByConsumerAndTenor(consumerID uint, tenor int) (*model.ConsumerLimit, error)
	GetConsumerLimits(consumerID uint) ([]model.ConsumerLimit, error)
	UpdateLimit(limit *model.ConsumerLimit) error
}

// TransactionUsecase defines all business logic operations for Transaction
type TransactionUsecase interface {
	CreateTransaction(transaction *model.Transaction) error
	GetTransaction(id uint) (*model.Transaction, error)
	GetConsumerTransactions(consumerID uint) ([]model.Transaction, error)
	UpdateTransactionStatus(id uint, status string) error
}

// consumerUsecase is the implementation of ConsumerUsecase
type consumerUsecase struct {
	repo repository.ConsumerRepository
	mu   sync.RWMutex
}

// NewConsumerUsecase creates a new instance of ConsumerUsecase
func NewConsumerUsecase(repo repository.ConsumerRepository) ConsumerUsecase {
	return &consumerUsecase{repo: repo}
}

// ValidateNIK validates Indonesian ID number format
func (u *consumerUsecase) ValidateNIK(nik string) error {
	if len(nik) != 16 {
		return errors.New("NIK must be 16 characters")
	}
	matched, err := regexp.MatchString("^[0-9]{16}$", nik)
	if err != nil || !matched {
		return errors.New("NIK must contain only numbers")
	}
	return nil
}

// RegisterConsumer registers a new consumer with validation
func (u *consumerUsecase) RegisterConsumer(consumer *model.Consumer) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Validation 1: Check if data is not empty
	if consumer.FullName == "" || consumer.NIK == "" || consumer.LegalName == "" {
		return errors.New("nama lengkap, nama sah, dan NIK tidak boleh kosong")
	}

	// Validation 2: Validate NIK format
	if err := u.ValidateNIK(consumer.NIK); err != nil {
		return err
	}

	// Validation 3: Check salary (Multifinance requirement)
	if consumer.Salary < 0 {
		return errors.New("gaji tidak valid (negatif)")
	}

	// Validation 4: Minimum salary requirement
	if consumer.Salary < 1000000 { // Minimum 1 juta
		return errors.New("gaji minimum 1 juta rupiah")
	}

	// Validation 5: Check date of birth
	if !consumer.DateOfBirth.IsZero() && consumer.DateOfBirth.After(time.Now()) {
		return errors.New("tanggal lahir tidak valid")
	}

	consumer.CreatedAt = time.Now()
	consumer.UpdatedAt = time.Now()

	log.Println("✓ Logika Bisnis OK. Menyimpan konsumen...")
	return u.repo.Create(consumer)
}

func (u *consumerUsecase) GetConsumer(id uint) (*model.Consumer, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.repo.GetByID(id)
}

func (u *consumerUsecase) GetConsumerByNIK(nik string) (*model.Consumer, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.repo.GetByNIK(nik)
}

func (u *consumerUsecase) UpdateConsumer(consumer *model.Consumer) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if consumer.ID == 0 {
		return errors.New("consumer ID tidak valid")
	}

	consumer.UpdatedAt = time.Now()
	return u.repo.Update(consumer)
}

func (u *consumerUsecase) DeleteConsumer(id uint) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if id == 0 {
		return errors.New("consumer ID tidak valid")
	}

	return u.repo.Delete(id)
}

// consumerLimitUsecase is the implementation of ConsumerLimitUsecase
type consumerLimitUsecase struct {
	limitRepo repository.ConsumerLimitRepository
	mu        sync.RWMutex
}

// NewConsumerLimitUsecase creates a new instance of ConsumerLimitUsecase
func NewConsumerLimitUsecase(limitRepo repository.ConsumerLimitRepository) ConsumerLimitUsecase {
	return &consumerLimitUsecase{limitRepo: limitRepo}
}

// AssignLimit assigns a credit limit to a consumer (ACID transaction)
func (u *consumerLimitUsecase) AssignLimit(limit *model.ConsumerLimit) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Validation 1: Check valid tenor
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[limit.Tenor] {
		return errors.New("tenor harus 1, 2, 3, atau 6 bulan")
	}

	// Validation 2: Check limit amount
	if limit.LimitAmount <= 0 {
		return errors.New("jumlah limit harus lebih dari 0")
	}

	// Validation 3: Check consumer exists
	if limit.ConsumerID == 0 {
		return errors.New("consumer ID tidak valid")
	}

	limit.UsedAmount = 0
	limit.CreatedAt = time.Now()
	limit.UpdatedAt = time.Now()

	log.Println("✓ Limit validation OK. Assigning limit...")
	return u.limitRepo.Create(limit)
}

func (u *consumerLimitUsecase) GetLimitByConsumerAndTenor(consumerID uint, tenor int) (*model.ConsumerLimit, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.limitRepo.GetByConsumerAndTenor(consumerID, tenor)
}

func (u *consumerLimitUsecase) GetConsumerLimits(consumerID uint) ([]model.ConsumerLimit, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()
	return u.limitRepo.GetByConsumerID(consumerID)
}

func (u *consumerLimitUsecase) UpdateLimit(limit *model.ConsumerLimit) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if limit.ID == 0 {
		return errors.New("limit ID tidak valid")
	}

	limit.UpdatedAt = time.Now()
	return u.limitRepo.Update(limit)
}

// transactionUsecase is the implementation of TransactionUsecase
type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	limitRepo       repository.ConsumerLimitRepository
	mu              sync.Mutex // Mutex for concurrent transaction handling
}

// NewTransactionUsecase creates a new instance of TransactionUsecase
func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	limitRepo repository.ConsumerLimitRepository,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		limitRepo:       limitRepo,
	}
}

// CreateTransaction creates a new transaction with concurrent limit checking
// This handles concurrent transactions as per requirement
func (u *transactionUsecase) CreateTransaction(transaction *model.Transaction) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	// Validation 1: Check required fields
	if transaction.ConsumerID == 0 || transaction.ContractNumber == "" {
		return errors.New("consumer ID dan nomor kontrak tidak boleh kosong")
	}

	// Validation 2: Check contract number uniqueness
	existingTx, err := u.transactionRepo.GetByContractNumber(transaction.ContractNumber)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	if existingTx != nil {
		return errors.New("nomor kontrak sudah digunakan")
	}

	// Validation 3: Check amount
	if transaction.OTR <= 0 {
		return errors.New("OTR harus lebih dari 0")
	}

	// Validation 4: Check tenor
	validTenors := map[int]bool{1: true, 2: true, 3: true, 6: true}
	if !validTenors[transaction.Tenor] {
		return errors.New("tenor harus 1, 2, 3, atau 6 bulan")
	}

	// CRITICAL: Check and deduct from consumer limit (ACID compliance)
	limit, err := u.limitRepo.GetByConsumerAndTenor(transaction.ConsumerID, transaction.Tenor)
	if err != nil {
		return errors.New("limit tidak ditemukan untuk tenor tersebut")
	}

	// Check if limit is sufficient
	if limit.UsedAmount+transaction.OTR > limit.LimitAmount {
		return errors.New("limit tidak cukup untuk transaksi ini")
	}

	// Deduct the limit (Atomicity - ensure this is transactional)
	limit.UsedAmount += transaction.OTR
	if err := u.limitRepo.Update(limit); err != nil {
		return errors.New("gagal update limit")
	}

	transaction.Status = "ACTIVE"
	transaction.CreatedAt = time.Now()
	transaction.UpdatedAt = time.Now()

	log.Println("✓ Transaction validation OK. Creating transaction...")
	return u.transactionRepo.Create(transaction)
}

func (u *transactionUsecase) GetTransaction(id uint) (*model.Transaction, error) {
	return u.transactionRepo.GetByID(id)
}

func (u *transactionUsecase) GetConsumerTransactions(consumerID uint) ([]model.Transaction, error) {
	return u.transactionRepo.GetByConsumerID(consumerID)
}

func (u *transactionUsecase) UpdateTransactionStatus(id uint, status string) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	validStatuses := map[string]bool{"ACTIVE": true, "COMPLETED": true, "DEFAULTED": true}
	if !validStatuses[status] {
		return errors.New("status tidak valid")
	}

	transaction, err := u.transactionRepo.GetByID(id)
	if err != nil {
		return err
	}

	transaction.Status = status
	transaction.UpdatedAt = time.Now()
	return u.transactionRepo.Update(transaction)
}
