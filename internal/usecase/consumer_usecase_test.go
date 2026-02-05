package usecase

import (
	"testing"
	"time"

	"main/internal/model"

	"gorm.io/gorm"
)

// MockConsumerRepository for testing
type MockConsumerRepository struct {
	consumers map[uint]*model.Consumer
	nextID    uint
}

func NewMockConsumerRepository() *MockConsumerRepository {
	return &MockConsumerRepository{
		consumers: make(map[uint]*model.Consumer),
		nextID:    1,
	}
}

func (m *MockConsumerRepository) Create(consumer *model.Consumer) error {
	consumer.ID = m.nextID
	m.consumers[m.nextID] = consumer
	m.nextID++
	return nil
}

func (m *MockConsumerRepository) GetByID(id uint) (*model.Consumer, error) {
	if consumer, exists := m.consumers[id]; exists {
		return consumer, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockConsumerRepository) GetByNIK(nik string) (*model.Consumer, error) {
	for _, consumer := range m.consumers {
		if consumer.NIK == nik {
			return consumer, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockConsumerRepository) GetAll() ([]model.Consumer, error) {
	var consumers []model.Consumer
	for _, consumer := range m.consumers {
		consumers = append(consumers, *consumer)
	}
	return consumers, nil
}

func (m *MockConsumerRepository) Update(consumer *model.Consumer) error {
	if _, exists := m.consumers[consumer.ID]; exists {
		m.consumers[consumer.ID] = consumer
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockConsumerRepository) Delete(id uint) error {
	delete(m.consumers, id)
	return nil
}

// Test: Valid Consumer Registration
func TestRegisterConsumer_Valid(t *testing.T) {
	mockRepo := NewMockConsumerRepository()
	uc := NewConsumerUsecase(mockRepo)

	consumer := &model.Consumer{
		NIK:         "1234567890123456",
		FullName:    "John Doe",
		LegalName:   "John Doe",
		Salary:      5000000,
		DateOfBirth: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}

	err := uc.RegisterConsumer(consumer)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if consumer.ID == 0 {
		t.Error("Expected consumer to have ID assigned")
	}
}

// Test: Invalid NIK Format
func TestRegisterConsumer_InvalidNIK(t *testing.T) {
	mockRepo := NewMockConsumerRepository()
	uc := NewConsumerUsecase(mockRepo)

	consumer := &model.Consumer{
		NIK:       "INVALID",
		FullName:  "John Doe",
		LegalName: "John Doe",
		Salary:    5000000,
	}

	err := uc.RegisterConsumer(consumer)
	if err == nil {
		t.Error("Expected error for invalid NIK, got nil")
	}
}

// Test: Missing Required Fields
func TestRegisterConsumer_MissingFields(t *testing.T) {
	mockRepo := NewMockConsumerRepository()
	uc := NewConsumerUsecase(mockRepo)

	consumer := &model.Consumer{
		NIK:      "1234567890123456",
		FullName: "",
		Salary:   5000000,
	}

	err := uc.RegisterConsumer(consumer)
	if err == nil {
		t.Error("Expected error for missing fields, got nil")
	}
}

// Test: Insufficient Salary
func TestRegisterConsumer_LowSalary(t *testing.T) {
	mockRepo := NewMockConsumerRepository()
	uc := NewConsumerUsecase(mockRepo)

	consumer := &model.Consumer{
		NIK:       "1234567890123456",
		FullName:  "John Doe",
		LegalName: "John Doe",
		Salary:    500000, // Below minimum
	}

	err := uc.RegisterConsumer(consumer)
	if err == nil {
		t.Error("Expected error for low salary, got nil")
	}
}

// Test: Get Consumer by ID
func TestGetConsumer(t *testing.T) {
	mockRepo := NewMockConsumerRepository()
	uc := NewConsumerUsecase(mockRepo)

	consumer := &model.Consumer{
		NIK:       "1234567890123456",
		FullName:  "John Doe",
		LegalName: "John Doe",
		Salary:    5000000,
	}
	uc.RegisterConsumer(consumer)

	retrieved, err := uc.GetConsumer(consumer.ID)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if retrieved.NIK != "1234567890123456" {
		t.Errorf("Expected NIK 1234567890123456, got %s", retrieved.NIK)
	}
}

// MockConsumerLimitRepository for testing
type MockConsumerLimitRepository struct {
	limits map[uint]*model.ConsumerLimit
	nextID uint
}

func NewMockConsumerLimitRepository() *MockConsumerLimitRepository {
	return &MockConsumerLimitRepository{
		limits: make(map[uint]*model.ConsumerLimit),
		nextID: 1,
	}
}

func (m *MockConsumerLimitRepository) Create(limit *model.ConsumerLimit) error {
	limit.ID = m.nextID
	m.limits[m.nextID] = limit
	m.nextID++
	return nil
}

func (m *MockConsumerLimitRepository) GetByID(id uint) (*model.ConsumerLimit, error) {
	if limit, exists := m.limits[id]; exists {
		return limit, nil
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockConsumerLimitRepository) GetByConsumerAndTenor(consumerID uint, tenor int) (*model.ConsumerLimit, error) {
	for _, limit := range m.limits {
		if limit.ConsumerID == consumerID && limit.Tenor == tenor {
			return limit, nil
		}
	}
	return nil, gorm.ErrRecordNotFound
}

func (m *MockConsumerLimitRepository) GetByConsumerID(consumerID uint) ([]model.ConsumerLimit, error) {
	var limits []model.ConsumerLimit
	for _, limit := range m.limits {
		if limit.ConsumerID == consumerID {
			limits = append(limits, *limit)
		}
	}
	return limits, nil
}

func (m *MockConsumerLimitRepository) Update(limit *model.ConsumerLimit) error {
	if _, exists := m.limits[limit.ID]; exists {
		m.limits[limit.ID] = limit
		return nil
	}
	return gorm.ErrRecordNotFound
}

func (m *MockConsumerLimitRepository) Delete(id uint) error {
	delete(m.limits, id)
	return nil
}

// Test: Valid Limit Assignment
func TestAssignLimit_Valid(t *testing.T) {
	mockRepo := NewMockConsumerLimitRepository()
	uc := NewConsumerLimitUsecase(mockRepo)

	limit := &model.ConsumerLimit{
		ConsumerID:  1,
		Tenor:       6,
		LimitAmount: 1000000,
	}

	err := uc.AssignLimit(limit)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if limit.UsedAmount != 0 {
		t.Errorf("Expected used amount to be 0, got %f", limit.UsedAmount)
	}
}

// Test: Invalid Tenor
func TestAssignLimit_InvalidTenor(t *testing.T) {
	mockRepo := NewMockConsumerLimitRepository()
	uc := NewConsumerLimitUsecase(mockRepo)

	limit := &model.ConsumerLimit{
		ConsumerID:  1,
		Tenor:       12, // Invalid tenor
		LimitAmount: 1000000,
	}

	err := uc.AssignLimit(limit)
	if err == nil {
		t.Error("Expected error for invalid tenor, got nil")
	}
}

// Test: Invalid Limit Amount
func TestAssignLimit_InvalidAmount(t *testing.T) {
	mockRepo := NewMockConsumerLimitRepository()
	uc := NewConsumerLimitUsecase(mockRepo)

	limit := &model.ConsumerLimit{
		ConsumerID:  1,
		Tenor:       6,
		LimitAmount: -1000000, // Negative amount
	}

	err := uc.AssignLimit(limit)
	if err == nil {
		t.Error("Expected error for negative limit, got nil")
	}
}
