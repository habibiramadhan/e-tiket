//test/usecase/transaction_usecase_test.go

package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ticket-system/internal/domain/entity"
	"ticket-system/internal/usecase"
	"ticket-system/test/mocks"
)

func TestCreateTransaction(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Transaction")).Return(1, nil).Once()
		mockEventRepo.On("UpdateTicketsSold", ctx, eventID, 2).Return(nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, 1, response.ID)
		assert.Equal(t, eventID, response.EventID)
		assert.Equal(t, event.Title, response.EventTitle)
		assert.Equal(t, req.Quantity, response.Quantity)
		assert.Equal(t, float64(req.Quantity)*event.Price, response.TotalAmount)
		assert.Equal(t, "pending", response.Status)
		assert.Equal(t, req.PaymentMethod, response.PaymentMethod)
		assert.NotEmpty(t, response.TransactionCode)
		assert.NotEmpty(t, response.PaymentDetail)
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("User Not Found", func(t *testing.T) {
		userID := 99
		
		req := usecase.CreateTransactionRequest{
			EventID:       1,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "pengguna tidak ditemukan", err.Error())
		
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found", func(t *testing.T) {
		userID := 1
		eventID := 99
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "event tidak ditemukan", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Active", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "cancelled",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "event tidak aktif", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Not Enough Tickets", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 999,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "jumlah tiket yang diminta melebihi kapasitas", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Quantity", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      0,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "jumlah tiket harus lebih dari 0", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Payment Method", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "invalid_method",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "metode pembayaran tidak valid", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error", func(t *testing.T) {
		userID := 1
		eventID := 1
		
		user := &entity.User{
			ID:       userID,
			Username: "testuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.CreateTransactionRequest{
			EventID:       eventID,
			Quantity:      2,
			PaymentMethod: "bank_transfer",
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		mockTransactionRepo.On("Create", ctx, mock.AnythingOfType("*entity.Transaction")).Return(0, errors.New("database error")).Once()
		
		response, err := transactionUsecase.CreateTransaction(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "database error", err.Error())
		
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetTransactionByID(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success - Owner", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		event := &entity.Event{
			ID:          2,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockEventRepo.On("FindByID", ctx, transaction.EventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.GetTransactionByID(ctx, userID, transactionID)
		
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, transactionID, response.ID)
		assert.Equal(t, transaction.TransactionCode, response.TransactionCode)
		assert.Equal(t, event.ID, response.EventID)
		assert.Equal(t, event.Title, response.EventTitle)
		
		mockTransactionRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Success - Organizer", func(t *testing.T) {
		userID := 2
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          1,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		organizer := &entity.User{
			ID:       userID,
			Username: "organizer",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		event := &entity.Event{
			ID:          2,
			Title:       "Konser Musik",
			Description: "Konser musik tahunan",
			Location:    "Jakarta Convention Center",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockUserRepo.On("FindByID", ctx, userID).Return(organizer, nil).Once()
		mockEventRepo.On("FindByID", ctx, transaction.EventID).Return(event, nil).Once()
		
		response, err := transactionUsecase.GetTransactionByID(ctx, userID, transactionID)
		
		assert.NoError(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, transactionID, response.ID)
		assert.Equal(t, transaction.TransactionCode, response.TransactionCode)
		assert.Equal(t, event.ID, response.EventID)
		assert.Equal(t, event.Title, response.EventTitle)
		
		mockTransactionRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Transaction Not Found", func(t *testing.T) {
		userID := 1
		transactionID := 99
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(nil, nil).Once()
		
		response, err := transactionUsecase.GetTransactionByID(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "transaksi tidak ditemukan", err.Error())
		
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Unauthorized Access", func(t *testing.T) {
		userID := 2
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          1,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		regularUser := &entity.User{
			ID:       userID,
			Username: "regularuser",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockUserRepo.On("FindByID", ctx, userID).Return(regularUser, nil).Once()
		
		response, err := transactionUsecase.GetTransactionByID(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "anda tidak memiliki izin untuk melihat transaksi ini", err.Error())
		
		mockTransactionRepo.AssertExpectations(t)
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         99,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockEventRepo.On("FindByID", ctx, transaction.EventID).Return(nil, nil).Once()
		
		response, err := transactionUsecase.GetTransactionByID(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Nil(t, response)
		assert.Equal(t, "event terkait tidak ditemukan", err.Error())
		
		mockTransactionRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
}

func TestUploadPaymentProof(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		userID := 1
		transactionID := 1
		proofURL := "https://example.com/proof.jpg"
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "1",
			ProofURL:      proofURL,
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockTransactionRepo.On("UpdatePaymentProof", ctx, transactionID, proofURL).Return(nil).Once()
		mockTransactionRepo.On("UpdateStatus", ctx, transactionID, "waiting_verification").Return(nil).Once()
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Transaction ID", func(t *testing.T) {
		userID := 1
		proofURL := "https://example.com/proof.jpg"
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "invalid",
			ProofURL:      proofURL,
		}
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "ID transaksi tidak valid", err.Error())
	})
	
	t.Run("Transaction Not Found", func(t *testing.T) {
		userID := 1
		transactionID := 99
		proofURL := "https://example.com/proof.jpg"
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "99",
			ProofURL:      proofURL,
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(nil, nil).Once()
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "transaksi tidak ditemukan", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Unauthorized Access", func(t *testing.T) {
		userID := 2
		transactionID := 1
		proofURL := "https://example.com/proof.jpg"
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          1,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "1",
			ProofURL:      proofURL,
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "anda tidak memiliki izin untuk transaksi ini", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Status", func(t *testing.T) {
		userID := 1
		transactionID := 1
		proofURL := "https://example.com/proof.jpg"
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "waiting_verification",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "1",
			ProofURL:      proofURL,
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "bukti pembayaran hanya dapat diunggah untuk transaksi dengan status pending", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Empty Proof URL", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		req := usecase.UploadPaymentProofRequest{
			TransactionID: "1",
			ProofURL:      "",
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.UploadPaymentProof(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "URL bukti pembayaran tidak boleh kosong", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestVerifyPayment(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		organizerID := 1
		transactionID := 1
		
		organizer := &entity.User{
			ID:       organizerID,
			Username: "organizer",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          2,
			EventID:         3,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "waiting_verification",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			PaymentProof:    "https://example.com/proof.jpg",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockUserRepo.On("FindByID", ctx, organizerID).Return(organizer, nil).Once()
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockTransactionRepo.On("VerifyPayment", ctx, transactionID, organizerID).Return(nil).Once()
		
		err := transactionUsecase.VerifyPayment(ctx, organizerID, transactionID)
		
		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Not an Organizer", func(t *testing.T) {
		regularUserID := 2
		transactionID := 1
		
		regularUser := &entity.User{
			ID:       regularUserID,
			Username: "user",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		mockUserRepo.On("FindByID", ctx, regularUserID).Return(regularUser, nil).Once()
		
		err := transactionUsecase.VerifyPayment(ctx, regularUserID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "hanya organizer yang dapat memverifikasi pembayaran", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Transaction Not Found", func(t *testing.T) {
		organizerID := 1
		transactionID := 99
		
		organizer := &entity.User{
			ID:       organizerID,
			Username: "organizer",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		mockUserRepo.On("FindByID", ctx, organizerID).Return(organizer, nil).Once()
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(nil, nil).Once()
		
		err := transactionUsecase.VerifyPayment(ctx, organizerID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "transaksi tidak ditemukan", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Status", func(t *testing.T) {
		organizerID := 1
		transactionID := 1
		
		organizer := &entity.User{
			ID:       organizerID,
			Username: "organizer",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          2,
			EventID:         3,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending", // Not waiting_verification
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			PaymentProof:    "https://example.com/proof.jpg",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockUserRepo.On("FindByID", ctx, organizerID).Return(organizer, nil).Once()
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.VerifyPayment(ctx, organizerID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "hanya transaksi dengan status menunggu verifikasi yang dapat diverifikasi", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestCancelTransaction(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockTransactionRepo.On("UpdateStatus", ctx, transactionID, "cancelled").Return(nil).Once()
		mockEventRepo.On("UpdateTicketsSold", ctx, transaction.EventID, -transaction.Quantity).Return(nil).Once()
		
		err := transactionUsecase.CancelTransaction(ctx, userID, transactionID)
		
		assert.NoError(t, err)
		mockTransactionRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Transaction Not Found", func(t *testing.T) {
		userID := 1
		transactionID := 99
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(nil, nil).Once()
		
		err := transactionUsecase.CancelTransaction(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "transaksi tidak ditemukan", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Unauthorized Access", func(t *testing.T) {
		userID := 2 // Not the owner
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          1, // Different from userID
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.CancelTransaction(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "anda tidak memiliki izin untuk transaksi ini", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Status", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "waiting_verification", // Not pending
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		
		err := transactionUsecase.CancelTransaction(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "hanya transaksi dengan status pending yang dapat dibatalkan", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error On Status Update", func(t *testing.T) {
		userID := 1
		transactionID := 1
		
		transaction := &entity.Transaction{
			ID:              transactionID,
			UserID:          userID,
			EventID:         2,
			TransactionCode: "TRX-20230101-123456",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Bank Transfer Details",
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		}
		
		mockTransactionRepo.On("FindByID", ctx, transactionID).Return(transaction, nil).Once()
		mockTransactionRepo.On("UpdateStatus", ctx, transactionID, "cancelled").Return(errors.New("database error")).Once()
		
		err := transactionUsecase.CancelTransaction(ctx, userID, transactionID)
		
		assert.Error(t, err)
		assert.Equal(t, "database error", err.Error())
		mockTransactionRepo.AssertExpectations(t)
	})
}

func TestGetUserTransactions(t *testing.T) {
	mockTransactionRepo := new(mocks.MockTransactionRepository)
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	transactionUsecase := usecase.NewTransactionUsecase(mockTransactionRepo, mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		userID := 1
		page := 1
		limit := 10
		offset := (page - 1) * limit
		
		transactions := []entity.Transaction{
			{
				ID:              1,
				UserID:          userID,
				EventID:         2,
				TransactionCode: "TRX-20230101-123456",
				Quantity:        2,
				TotalAmount:     500000,
				Status:          "pending",
				PaymentMethod:   "bank_transfer",
				PaymentDetail:   "Bank Transfer Details",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              2,
				UserID:          userID,
				EventID:         3,
				TransactionCode: "TRX-20230102-654321",
				Quantity:        1,
				TotalAmount:     250000,
				Status:          "success",
				PaymentMethod:   "qris",
				PaymentDetail:   "QRIS Payment Details",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		
		event1 := &entity.Event{
			ID:    2,
			Title: "Konser Musik",
		}
		
		event2 := &entity.Event{
			ID:    3,
			Title: "Festival Film",
		}
		
		total := 2
		
		mockTransactionRepo.On("FindByUserID", ctx, userID, offset, limit).Return(transactions, nil).Once()
		mockTransactionRepo.On("CountByUserID", ctx, userID).Return(total, nil).Once()
		mockEventRepo.On("FindByID", ctx, transactions[0].EventID).Return(event1, nil).Once()
		mockEventRepo.On("FindByID", ctx, transactions[1].EventID).Return(event2, nil).Once()
		
		responses, count, err := transactionUsecase.GetUserTransactions(ctx, userID, page, limit)
		
		assert.NoError(t, err)
		assert.Equal(t, total, count)
		assert.Equal(t, len(transactions), len(responses))
		
		assert.Equal(t, transactions[0].ID, responses[0].ID)
		assert.Equal(t, transactions[0].TransactionCode, responses[0].TransactionCode)
		assert.Equal(t, event1.Title, responses[0].EventTitle)
		
		assert.Equal(t, transactions[1].ID, responses[1].ID)
		assert.Equal(t, transactions[1].TransactionCode, responses[1].TransactionCode)
		assert.Equal(t, event2.Title, responses[1].EventTitle)
		
		mockTransactionRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Empty Result", func(t *testing.T) {
		userID := 1
		page := 1
		limit := 10
		offset := (page - 1) * limit
		
		var transactions []entity.Transaction
		total := 0
		
		mockTransactionRepo.On("FindByUserID", ctx, userID, offset, limit).Return(transactions, nil).Once()
		mockTransactionRepo.On("CountByUserID", ctx, userID).Return(total, nil).Once()
		
		responses, count, err := transactionUsecase.GetUserTransactions(ctx, userID, page, limit)
		
		assert.NoError(t, err)
		assert.Equal(t, total, count)
		assert.Equal(t, 0, len(responses))
		
		mockTransactionRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found Error Handling", func(t *testing.T) {
		userID := 1
		page := 1
		limit := 10
		offset := (page - 1) * limit
		
		transactions := []entity.Transaction{
			{
				ID:              1,
				UserID:          userID,
				EventID:         2,
				TransactionCode: "TRX-20230101-123456",
				Quantity:        2,
				TotalAmount:     500000,
				Status:          "pending",
				PaymentMethod:   "bank_transfer",
				PaymentDetail:   "Bank Transfer Details",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		
		total := 1
		
		mockTransactionRepo.On("FindByUserID", ctx, userID, offset, limit).Return(transactions, nil).Once()
		mockTransactionRepo.On("CountByUserID", ctx, userID).Return(total, nil).Once()
		mockEventRepo.On("FindByID", ctx, transactions[0].EventID).Return(nil, nil).Once()
		
		responses, count, err := transactionUsecase.GetUserTransactions(ctx, userID, page, limit)
		
		assert.NoError(t, err)
		assert.Equal(t, total, count)
		assert.Equal(t, len(transactions), len(responses))
		assert.Equal(t, "Event Tidak Ditemukan", responses[0].EventTitle)
		
		mockTransactionRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error", func(t *testing.T) {
		userID := 1
		page := 1
		limit := 10
		offset := (page - 1) * limit
		
		mockTransactionRepo.On("FindByUserID", ctx, userID, offset, limit).Return(nil, errors.New("database error")).Once()
		
		responses, count, err := transactionUsecase.GetUserTransactions(ctx, userID, page, limit)
		
		assert.Error(t, err)
		assert.Equal(t, 0, count)
		assert.Nil(t, responses)
		assert.Equal(t, "database error", err.Error())
		
		mockTransactionRepo.AssertExpectations(t)
	})
}