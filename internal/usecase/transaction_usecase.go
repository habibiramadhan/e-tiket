//internal/usecase/transaction_usecase.go

package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	
	"ticket-system/internal/domain/entity"
	"ticket-system/internal/domain/repository"
	"ticket-system/pkg/utils"
)

type CreateTransactionRequest struct {
	EventID       int    `json:"event_id"`
	Quantity      int    `json:"quantity"`
	PaymentMethod string `json:"payment_method"`
}

type TransactionResponse struct {
	ID              int       `json:"id"`
	TransactionCode string    `json:"transaction_code"`
	EventID         int       `json:"event_id"`
	EventTitle      string    `json:"event_title"`
	Quantity        int       `json:"quantity"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentDetail   string    `json:"payment_detail"`
	PaymentProof    string    `json:"payment_proof,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
}

type UploadPaymentProofRequest struct {
	TransactionID string `json:"transaction_id"`
	ProofURL      string `json:"proof_url"`
}

type TransactionUsecase interface {
	CreateTransaction(ctx context.Context, userID int, req CreateTransactionRequest) (*TransactionResponse, error)
	GetTransactionByID(ctx context.Context, userID int, transactionID int) (*TransactionResponse, error)
	GetTransactionByCode(ctx context.Context, userID int, code string) (*TransactionResponse, error)
	GetUserTransactions(ctx context.Context, userID, page, limit int) ([]TransactionResponse, int, error)
	UploadPaymentProof(ctx context.Context, userID int, req UploadPaymentProofRequest) error
	CancelTransaction(ctx context.Context, userID int, transactionID int) error
	VerifyPayment(ctx context.Context, organizerID int, transactionID int) error
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	eventRepo       repository.EventRepository
	userRepo        repository.UserRepository
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	eventRepo repository.EventRepository,
	userRepo repository.UserRepository,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		eventRepo:       eventRepo,
		userRepo:        userRepo,
	}
}

func (u *transactionUsecase) CreateTransaction(ctx context.Context, userID int, req CreateTransactionRequest) (*TransactionResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("pengguna tidak ditemukan")
	}

	event, err := u.eventRepo.FindByID(ctx, req.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event tidak ditemukan")
	}

	if event.Status != "active" {
		return nil, errors.New("event tidak aktif")
	}

	if event.TicketsSold + req.Quantity > event.MaxCapacity {
		return nil, errors.New("jumlah tiket yang diminta melebihi kapasitas")
	}

	if req.Quantity <= 0 {
		return nil, errors.New("jumlah tiket harus lebih dari 0")
	}

	if req.PaymentMethod == "" {
		return nil, errors.New("metode pembayaran harus dipilih")
	}

	validPaymentMethods := map[string]bool{
		"bank_transfer": true,
		"qris":          true,
		"ewallet":       true,
	}

	if !validPaymentMethods[req.PaymentMethod] {
		return nil, errors.New("metode pembayaran tidak valid")
	}

	totalAmount := float64(req.Quantity) * event.Price

	transactionCode := fmt.Sprintf("TRX-%s-%d", time.Now().Format("20060102"), utils.GenerateRandomNumber(6))

	var paymentDetail string
	switch req.PaymentMethod {
	case "bank_transfer":
		paymentDetail = "Silakan transfer ke Bank BCA 1234567890 a/n Ticket System"
	case "qris":
		paymentDetail = "Silakan scan QRIS yang tersedia"
	case "ewallet":
		paymentDetail = "Silakan bayar melalui e-wallet yang terdaftar"
	}

	transaction := &entity.Transaction{
		UserID:          userID,
		EventID:         req.EventID,
		TransactionCode: transactionCode,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		PaymentDetail:   paymentDetail,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	transactionID, err := u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	err = u.eventRepo.UpdateTicketsSold(ctx, req.EventID, req.Quantity)
	if err != nil {
		return nil, err
	}

	response := &TransactionResponse{
		ID:              transactionID,
		TransactionCode: transactionCode,
		EventID:         event.ID,
		EventTitle:      event.Title,
		Quantity:        req.Quantity,
		TotalAmount:     totalAmount,
		Status:          "pending",
		PaymentMethod:   req.PaymentMethod,
		PaymentDetail:   paymentDetail,
		CreatedAt:       transaction.CreatedAt,
	}

	return response, nil
}

func (u *transactionUsecase) GetTransactionByID(ctx context.Context, userID int, transactionID int) (*TransactionResponse, error) {
	transaction, err := u.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaction.UserID != userID {
		user, err := u.userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user == nil || user.Role != "organizer" {
			return nil, errors.New("anda tidak memiliki izin untuk melihat transaksi ini")
		}
	}

	event, err := u.eventRepo.FindByID(ctx, transaction.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event terkait tidak ditemukan")
	}

	response := &TransactionResponse{
		ID:              transaction.ID,
		TransactionCode: transaction.TransactionCode,
		EventID:         event.ID,
		EventTitle:      event.Title,
		Quantity:        transaction.Quantity,
		TotalAmount:     transaction.TotalAmount,
		Status:          transaction.Status,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentDetail:   transaction.PaymentDetail,
		PaymentProof:    transaction.PaymentProof,
		CreatedAt:       transaction.CreatedAt,
	}

	return response, nil
}

func (u *transactionUsecase) GetTransactionByCode(ctx context.Context, userID int, code string) (*TransactionResponse, error) {
	transaction, err := u.transactionRepo.FindByCode(ctx, code)
	if err != nil {
		return nil, err
	}

	if transaction == nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if transaction.UserID != userID {
		user, err := u.userRepo.FindByID(ctx, userID)
		if err != nil {
			return nil, err
		}
		if user == nil || user.Role != "organizer" {
			return nil, errors.New("anda tidak memiliki izin untuk melihat transaksi ini")
		}
	}

	event, err := u.eventRepo.FindByID(ctx, transaction.EventID)
	if err != nil {
		return nil, err
	}
	if event == nil {
		return nil, errors.New("event terkait tidak ditemukan")
	}

	response := &TransactionResponse{
		ID:              transaction.ID,
		TransactionCode: transaction.TransactionCode,
		EventID:         event.ID,
		EventTitle:      event.Title,
		Quantity:        transaction.Quantity,
		TotalAmount:     transaction.TotalAmount,
		Status:          transaction.Status,
		PaymentMethod:   transaction.PaymentMethod,
		PaymentDetail:   transaction.PaymentDetail,
		PaymentProof:    transaction.PaymentProof,
		CreatedAt:       transaction.CreatedAt,
	}

	return response, nil
}

func (u *transactionUsecase) GetUserTransactions(ctx context.Context, userID, page, limit int) ([]TransactionResponse, int, error) {
	offset := (page - 1) * limit
	transactions, err := u.transactionRepo.FindByUserID(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}

	total, err := u.transactionRepo.CountByUserID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}

	var responses []TransactionResponse
	for _, transaction := range transactions {
		event, err := u.eventRepo.FindByID(ctx, transaction.EventID)
		if err != nil {
			return nil, 0, err
		}
		
		eventTitle := "Event Tidak Ditemukan"
		if event != nil {
			eventTitle = event.Title
		}

		response := TransactionResponse{
			ID:              transaction.ID,
			TransactionCode: transaction.TransactionCode,
			EventID:         transaction.EventID,
			EventTitle:      eventTitle,
			Quantity:        transaction.Quantity,
			TotalAmount:     transaction.TotalAmount,
			Status:          transaction.Status,
			PaymentMethod:   transaction.PaymentMethod,
			PaymentDetail:   transaction.PaymentDetail,
			PaymentProof:    transaction.PaymentProof,
			CreatedAt:       transaction.CreatedAt,
		}

		responses = append(responses, response)
	}

	return responses, total, nil
}

func (u *transactionUsecase) UploadPaymentProof(ctx context.Context, userID int, req UploadPaymentProofRequest) error {
	transactionID, ok := utils.ParseStringToInt(req.TransactionID)
	if !ok {
		return errors.New("ID transaksi tidak valid")
	}
	
	transaction, err := u.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("transaksi tidak ditemukan")
	}

	if transaction.UserID != userID {
		return errors.New("anda tidak memiliki izin untuk transaksi ini")
	}

	if transaction.Status != "pending" {
		return errors.New("bukti pembayaran hanya dapat diunggah untuk transaksi dengan status pending")
	}

	if req.ProofURL == "" {
		return errors.New("URL bukti pembayaran tidak boleh kosong")
	}

	err = u.transactionRepo.UpdatePaymentProof(ctx, transactionID, req.ProofURL)
	if err != nil {
		return err
	}

	return u.transactionRepo.UpdateStatus(ctx, transactionID, "waiting_verification")
}

func (u *transactionUsecase) CancelTransaction(ctx context.Context, userID int, transactionID int) error {
	transaction, err := u.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("transaksi tidak ditemukan")
	}

	if transaction.UserID != userID {
		return errors.New("anda tidak memiliki izin untuk transaksi ini")
	}

	if transaction.Status != "pending" {
		return errors.New("hanya transaksi dengan status pending yang dapat dibatalkan")
	}

	err = u.transactionRepo.UpdateStatus(ctx, transactionID, "cancelled")
	if err != nil {
		return err
	}

	return u.eventRepo.UpdateTicketsSold(ctx, transaction.EventID, -transaction.Quantity)
}

func (u *transactionUsecase) VerifyPayment(ctx context.Context, organizerID int, transactionID int) error {
	organizer, err := u.userRepo.FindByID(ctx, organizerID)
	if err != nil {
		return err
	}

	if organizer == nil || organizer.Role != "organizer" {
		return errors.New("hanya organizer yang dapat memverifikasi pembayaran")
	}

	transaction, err := u.transactionRepo.FindByID(ctx, transactionID)
	if err != nil {
		return err
	}

	if transaction == nil {
		return errors.New("transaksi tidak ditemukan")
	}

	if transaction.Status != "waiting_verification" {
		return errors.New("hanya transaksi dengan status menunggu verifikasi yang dapat diverifikasi")
	}

	return u.transactionRepo.VerifyPayment(ctx, transactionID, organizerID)
}