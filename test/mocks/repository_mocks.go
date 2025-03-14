//test/mocks/repository_mocks.go

package mocks

import (
	"context"
	
	"github.com/stretchr/testify/mock"
	
	"ticket-system/internal/domain/entity"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(ctx context.Context, user *entity.User) (int, error) {
	args := m.Called(ctx, user)
	return args.Int(0), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, id int) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*entity.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, user *entity.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateVerificationStatus(ctx context.Context, userID int, isVerified bool) error {
	args := m.Called(ctx, userID, isVerified)
	return args.Error(0)
}

func (m *MockUserRepository) CreateDefaultProfile(ctx context.Context, userID int) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

type MockEventRepository struct {
	mock.Mock
}

func (m *MockEventRepository) Create(ctx context.Context, event *entity.Event) (int, error) {
	args := m.Called(ctx, event)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) FindByID(ctx context.Context, id int) (*entity.Event, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Event), args.Error(1)
}

func (m *MockEventRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.Event, error) {
	args := m.Called(ctx, offset, limit)
	return args.Get(0).([]entity.Event), args.Error(1)
}

func (m *MockEventRepository) CountAll(ctx context.Context) (int, error) {
	args := m.Called(ctx)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) FindByOwnerID(ctx context.Context, ownerID, offset, limit int) ([]entity.Event, error) {
	args := m.Called(ctx, ownerID, offset, limit)
	return args.Get(0).([]entity.Event), args.Error(1)
}

func (m *MockEventRepository) CountByOwnerID(ctx context.Context, ownerID int) (int, error) {
	args := m.Called(ctx, ownerID)
	return args.Int(0), args.Error(1)
}

func (m *MockEventRepository) Update(ctx context.Context, event *entity.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockEventRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEventRepository) UpdateTicketsSold(ctx context.Context, eventID, quantity int) error {
	args := m.Called(ctx, eventID, quantity)
	return args.Error(0)
}

type MockTransactionRepository struct {
	mock.Mock
}

func (m *MockTransactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (int, error) {
	args := m.Called(ctx, transaction)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) FindByID(ctx context.Context, id int) (*entity.Transaction, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByCode(ctx context.Context, code string) (*entity.Transaction, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) FindByUserID(ctx context.Context, userID, offset, limit int) ([]entity.Transaction, error) {
	args := m.Called(ctx, userID, offset, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]entity.Transaction), args.Error(1)
}

func (m *MockTransactionRepository) CountByUserID(ctx context.Context, userID int) (int, error) {
	args := m.Called(ctx, userID)
	return args.Int(0), args.Error(1)
}

func (m *MockTransactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	args := m.Called(ctx, transaction)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockTransactionRepository) UpdatePaymentProof(ctx context.Context, id int, proofURL string) error {
	args := m.Called(ctx, id, proofURL)
	return args.Error(0)
}

func (m *MockTransactionRepository) VerifyPayment(ctx context.Context, id, verifierID int) error {
	args := m.Called(ctx, id, verifierID)
	return args.Error(0)
}