//test/handler/transaction_handler_test.go

package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ticket-system/internal/delivery/http/handler"
	"ticket-system/internal/domain/entity"
	"ticket-system/internal/usecase"
	"ticket-system/pkg/utils"
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

type MockUserProfileRepository struct {
	mock.Mock
}

func (m *MockUserProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) (int, error) {
	args := m.Called(ctx, profile)
	return args.Int(0), args.Error(1)
}

func (m *MockUserProfileRepository) FindByUserID(ctx context.Context, userID int) (*entity.UserProfile, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.UserProfile), args.Error(1)
}

func (m *MockUserProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	args := m.Called(ctx, profile)
	return args.Error(0)
}

func (m *MockUserProfileRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

type MockEmailVerificationRepository struct {
	mock.Mock
}

func (m *MockEmailVerificationRepository) Create(ctx context.Context, verification *entity.EmailVerification) (int, error) {
	args := m.Called(ctx, verification)
	return args.Int(0), args.Error(1)
}

func (m *MockEmailVerificationRepository) FindByToken(ctx context.Context, token string) (*entity.EmailVerification, error) {
	args := m.Called(ctx, token)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.EmailVerification), args.Error(1)
}

func (m *MockEmailVerificationRepository) Delete(ctx context.Context, id int) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockEmailVerificationRepository) DeleteByUserID(ctx context.Context, userID int) error {
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

type MockTransactionUsecase struct {
	mock.Mock
}

func (m *MockTransactionUsecase) CreateTransaction(ctx context.Context, userID int, req usecase.CreateTransactionRequest) (*usecase.TransactionResponse, error) {
	args := m.Called(ctx, userID, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TransactionResponse), args.Error(1)
}

func (m *MockTransactionUsecase) GetTransactionByID(ctx context.Context, userID, transactionID int) (*usecase.TransactionResponse, error) {
	args := m.Called(ctx, userID, transactionID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TransactionResponse), args.Error(1)
}

func (m *MockTransactionUsecase) GetTransactionByCode(ctx context.Context, userID int, code string) (*usecase.TransactionResponse, error) {
	args := m.Called(ctx, userID, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.TransactionResponse), args.Error(1)
}

func (m *MockTransactionUsecase) GetUserTransactions(ctx context.Context, userID, page, limit int) ([]usecase.TransactionResponse, int, error) {
	args := m.Called(ctx, userID, page, limit)
	return args.Get(0).([]usecase.TransactionResponse), args.Int(1), args.Error(2)
}

func (m *MockTransactionUsecase) UploadPaymentProof(ctx context.Context, userID int, req usecase.UploadPaymentProofRequest) error {
	args := m.Called(ctx, userID, req)
	return args.Error(0)
}

func (m *MockTransactionUsecase) CancelTransaction(ctx context.Context, userID, transactionID int) error {
	args := m.Called(ctx, userID, transactionID)
	return args.Error(0)
}

func (m *MockTransactionUsecase) VerifyPayment(ctx context.Context, organizerID, transactionID int) error {
	args := m.Called(ctx, organizerID, transactionID)
	return args.Error(0)
}

type MockUserUsecase struct {
	mock.Mock
}

func (m *MockUserUsecase) Register(ctx context.Context, req usecase.RegisterRequest) (int, error) {
	args := m.Called(ctx, req)
	return args.Int(0), args.Error(1)
}

func (m *MockUserUsecase) Login(ctx context.Context, req usecase.LoginRequest) (*usecase.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*usecase.LoginResponse), args.Error(1)
}

func (m *MockUserUsecase) GetByID(ctx context.Context, id int) (*entity.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserUsecase) VerifyEmail(ctx context.Context, token string) error {
	args := m.Called(ctx, token)
	return args.Error(0)
}

func (m *MockUserUsecase) ResendVerificationEmail(ctx context.Context, email string) error {
	args := m.Called(ctx, email)
	return args.Error(0)
}

func (m *MockUserUsecase) UpdateProfile(ctx context.Context, userID int, name, gender, address, phoneNumber string) error {
	args := m.Called(ctx, userID, name, gender, address, phoneNumber)
	return args.Error(0)
}

func setupUserHandlerTest() (*fiber.App, *MockUserUsecase) {
	mockUserUsecase := new(MockUserUsecase)
	app := fiber.New()
	
	userHandler := handler.NewUserHandler(mockUserUsecase)
	
	app.Post("/api/register", userHandler.Register)
	app.Post("/api/login", userHandler.Login)
	app.Get("/api/verify-email", userHandler.VerifyEmail)
	
	return app, mockUserUsecase
}

func setupTransactionHandlerTest() (*fiber.App, *MockTransactionUsecase) {
	mockUsecase := new(MockTransactionUsecase)
	app := fiber.New()
	
	transactionHandler := handler.NewTransactionHandler(mockUsecase)
	
	app.Post("/api/transactions", func(c *fiber.Ctx) error {
		c.Locals("claims", &utils.JWTClaim{
			UserID:   1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		})
		return transactionHandler.CreateTransaction(c)
	})
	
	app.Get("/api/transactions/:id", func(c *fiber.Ctx) error {
		c.Locals("claims", &utils.JWTClaim{
			UserID:   1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		})
		return transactionHandler.GetTransactionByID(c)
	})
	
	app.Post("/api/transactions/proof", func(c *fiber.Ctx) error {
		c.Locals("claims", &utils.JWTClaim{
			UserID:   1,
			Username: "testuser",
			Email:    "test@example.com",
			Role:     "user",
		})
		return transactionHandler.UploadPaymentProof(c)
	})
	
	app.Put("/api/organizer/transactions/:id/verify", func(c *fiber.Ctx) error {
		c.Locals("claims", &utils.JWTClaim{
			UserID:   2,
			Username: "organizer",
			Email:    "organizer@example.com",
			Role:     "organizer",
		})
		return transactionHandler.VerifyPayment(c)
	})
	
	return app, mockUsecase
}

func TestUserFlow(t *testing.T) {
	app, mockUsecase := setupUserHandlerTest()
	
	t.Run("Register User", func(t *testing.T) {
		mockUsecase.On("Register", mock.Anything, mock.AnythingOfType("usecase.RegisterRequest")).Return(1, nil)
		
		reqBody := map[string]interface{}{
			"username":        "testuser",
			"email":           "testuser@example.com",
			"password":        "password123",
			"retype_password": "password123",
			"role":            "user",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest("POST", "/api/register", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
	})
	
	t.Run("Login User", func(t *testing.T) {
		loginResp := &usecase.LoginResponse{
			Token:        "dummy_token_123456",
			RefreshToken: "dummy_refresh_token_123456",
			Role:         "user",
			Username:     "testuser",
		}
		
		mockUsecase.On("Login", mock.Anything, mock.AnythingOfType("usecase.LoginRequest")).Return(loginResp, nil)
		
		reqBody := map[string]interface{}{
			"username": "testuser",
			"password": "password123",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest("POST", "/api/login", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
		data := result["data"].(map[string]interface{})
		assert.NotEmpty(t, data["token"])
	})
}

func TestTransactionFlow(t *testing.T) {
	app, mockUsecase := setupTransactionHandlerTest()
	
	var transactionID int
	t.Run("Create Transaction", func(t *testing.T) {
		mockResponse := &usecase.TransactionResponse{
			ID:              1,
			TransactionCode: "TRX-20230101-123456",
			EventID:         1,
			EventTitle:      "Konser Musik",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "pending",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Silakan transfer ke Bank BCA 1234567890 a/n Ticket System",
			CreatedAt:       time.Now(),
		}
		
		mockUsecase.On("CreateTransaction", mock.Anything, mock.Anything, mock.AnythingOfType("usecase.CreateTransactionRequest")).Return(mockResponse, nil)
		
		reqBody := map[string]interface{}{
			"event_id":       1,
			"quantity":       2,
			"payment_method": "bank_transfer",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest("POST", "/api/transactions", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
		
		data := result["data"].(map[string]interface{})
		transactionID = int(data["id"].(float64))
	})
	
	t.Run("Upload Payment Proof", func(t *testing.T) {
		mockUsecase.On("UploadPaymentProof", mock.Anything, mock.Anything, mock.AnythingOfType("usecase.UploadPaymentProofRequest")).Return(nil)
		
		reqBody := map[string]interface{}{
			"transaction_id": strconv.Itoa(transactionID),
			"proof_url":      "https://example.com/bukti-transfer.jpg",
		}
		
		jsonBody, _ := json.Marshal(reqBody)
		
		req, _ := http.NewRequest("POST", "/api/transactions/proof", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
	})
	
	t.Run("Verify Payment", func(t *testing.T) {
		mockUsecase.On("VerifyPayment", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		
		req, _ := http.NewRequest("PUT", "/api/organizer/transactions/"+strconv.Itoa(transactionID)+"/verify", nil)
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
	})
	
	t.Run("Get Transaction Detail", func(t *testing.T) {
		mockResponse := &usecase.TransactionResponse{
			ID:              transactionID,
			TransactionCode: "TRX-20230101-123456",
			EventID:         1,
			EventTitle:      "Konser Musik",
			Quantity:        2,
			TotalAmount:     500000,
			Status:          "success",
			PaymentMethod:   "bank_transfer",
			PaymentDetail:   "Silakan transfer ke Bank BCA 1234567890 a/n Ticket System",
			PaymentProof:    "https://example.com/bukti-transfer.jpg",
			CreatedAt:       time.Now(),
		}
		
		mockUsecase.On("GetTransactionByID", mock.Anything, mock.Anything, mock.Anything).Return(mockResponse, nil)
		
		req, _ := http.NewRequest("GET", "/api/transactions/"+strconv.Itoa(transactionID), nil)
		
		resp, err := app.Test(req)
		assert.NoError(t, err)
		
		bodyBytes, _ := io.ReadAll(resp.Body)
		
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		
		var result map[string]interface{}
		json.Unmarshal(bodyBytes, &result)
		
		assert.Equal(t, true, result["status"])
		
		data := result["data"].(map[string]interface{})
		assert.Equal(t, "success", data["status"])
	})
}