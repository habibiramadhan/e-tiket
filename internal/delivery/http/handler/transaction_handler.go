//internal/delivery/http/handler/transaction_handler.go

package handler

import (
	"log"
	"strconv"
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/usecase"
	"ticket-system/pkg/utils"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	log.Println("CreateTransaction handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	var req usecase.CreateTransactionRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	var validationErrors []utils.ErrorDetail
	
	if req.EventID <= 0 {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "event_id",
			Message: "ID event tidak valid",
		})
	}
	
	if req.Quantity <= 0 {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "quantity",
			Message: "Jumlah tiket harus lebih dari 0",
		})
	}
	
	if req.PaymentMethod == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "payment_method",
			Message: "Metode pembayaran tidak boleh kosong",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	response, err := h.transactionUsecase.CreateTransaction(c.Context(), userID, req)
	if err != nil {
		switch err.Error() {
		case "pengguna tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Pengguna tidak ditemukan", fiber.StatusNotFound)
		case "event tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event tidak ditemukan", fiber.StatusNotFound)
		case "event tidak aktif":
			return utils.ErrorResponse(c, utils.ErrorCodeEventCancelled, "Event tidak aktif", fiber.StatusBadRequest)
		case "jumlah tiket yang diminta melebihi kapasitas":
			return utils.ErrorResponse(c, utils.ErrorCodeTicketSoldOut, "Tidak cukup tiket tersedia", fiber.StatusBadRequest)
		case "jumlah tiket harus lebih dari 0":
			return utils.ErrorResponse(c, utils.ErrorCodeTicketInvalidQuantity, "Jumlah tiket harus lebih dari 0", fiber.StatusBadRequest)
		case "metode pembayaran harus dipilih":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Metode pembayaran harus dipilih", fiber.StatusBadRequest)
		case "metode pembayaran tidak valid":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Metode pembayaran tidak valid", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal membuat transaksi: "+err.Error())
		}
	}
	
	return utils.CreatedResponse(c, "Transaksi berhasil dibuat", response)
}

func (h *TransactionHandler) GetTransactionByID(c *fiber.Ctx) error {
	log.Println("GetTransactionByID handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID transaksi tidak valid", fiber.StatusBadRequest)
	}
	
	transaction, err := h.transactionUsecase.GetTransactionByID(c.Context(), userID, transactionID)
	if err != nil {
		switch err.Error() {
		case "transaksi tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Transaksi tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk melihat transaksi ini":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk melihat transaksi ini", fiber.StatusForbidden)
		case "event terkait tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event terkait tidak ditemukan", fiber.StatusNotFound)
		default:
			return utils.ServerError(c, "Gagal mendapatkan detail transaksi: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Detail transaksi berhasil diambil", transaction)
}

func (h *TransactionHandler) GetTransactionByCode(c *fiber.Ctx) error {
	log.Println("GetTransactionByCode handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	code := c.Query("code")
	if code == "" {
		return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "Kode transaksi diperlukan", fiber.StatusBadRequest)
	}
	
	transaction, err := h.transactionUsecase.GetTransactionByCode(c.Context(), userID, code)
	if err != nil {
		switch err.Error() {
		case "transaksi tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Transaksi tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk melihat transaksi ini":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk melihat transaksi ini", fiber.StatusForbidden)
		case "event terkait tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event terkait tidak ditemukan", fiber.StatusNotFound)
		default:
			return utils.ServerError(c, "Gagal mendapatkan detail transaksi: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Detail transaksi berhasil diambil", transaction)
}

func (h *TransactionHandler) GetUserTransactions(c *fiber.Ctx) error {
	log.Println("GetUserTransactions handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	transactions, total, err := h.transactionUsecase.GetUserTransactions(c.Context(), userID, page, limit)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan daftar transaksi: "+err.Error())
	}
	
	meta := fiber.Map{
		"page":  page,
		"limit": limit,
		"total": total,
	}
	
	return utils.SuccessResponse(c, "Daftar transaksi berhasil diambil", transactions, meta)
}

func (h *TransactionHandler) UploadPaymentProof(c *fiber.Ctx) error {
	log.Println("UploadPaymentProof handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	var req usecase.UploadPaymentProofRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	var validationErrors []utils.ErrorDetail
	
	if req.TransactionID == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "transaction_id",
			Message: "ID transaksi tidak boleh kosong",
		})
	}
	
	if req.ProofURL == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "proof_url",
			Message: "URL bukti pembayaran tidak boleh kosong",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	err = h.transactionUsecase.UploadPaymentProof(c.Context(), userID, req)
	if err != nil {
		switch err.Error() {
		case "ID transaksi tidak valid":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID transaksi tidak valid", fiber.StatusBadRequest)
		case "transaksi tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Transaksi tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk transaksi ini":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk transaksi ini", fiber.StatusForbidden)
		case "bukti pembayaran hanya dapat diunggah untuk transaksi dengan status pending":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Bukti pembayaran hanya dapat diunggah untuk transaksi dengan status pending", fiber.StatusBadRequest)
		case "URL bukti pembayaran tidak boleh kosong":
			return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "URL bukti pembayaran tidak boleh kosong", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal mengunggah bukti pembayaran: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Bukti pembayaran berhasil diunggah", nil)
}

func (h *TransactionHandler) CancelTransaction(c *fiber.Ctx) error {
	log.Println("CancelTransaction handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID transaksi tidak valid", fiber.StatusBadRequest)
	}
	
	err = h.transactionUsecase.CancelTransaction(c.Context(), userID, transactionID)
	if err != nil {
		switch err.Error() {
		case "transaksi tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Transaksi tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk transaksi ini":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk transaksi ini", fiber.StatusForbidden)
		case "hanya transaksi dengan status pending yang dapat dibatalkan":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Hanya transaksi dengan status pending yang dapat dibatalkan", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal membatalkan transaksi: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Transaksi berhasil dibatalkan", nil)
}

func (h *TransactionHandler) VerifyPayment(c *fiber.Ctx) error {
	log.Println("VerifyPayment handler called with path:", c.Path())
	
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	organizerID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	transactionID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID transaksi tidak valid", fiber.StatusBadRequest)
	}
	
	err = h.transactionUsecase.VerifyPayment(c.Context(), organizerID, transactionID)
	if err != nil {
		switch err.Error() {
		case "hanya organizer yang dapat memverifikasi pembayaran":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Hanya organizer yang dapat memverifikasi pembayaran", fiber.StatusForbidden)
		case "transaksi tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Transaksi tidak ditemukan", fiber.StatusNotFound)
		case "hanya transaksi dengan status menunggu verifikasi yang dapat diverifikasi":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Hanya transaksi dengan status menunggu verifikasi yang dapat diverifikasi", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal memverifikasi pembayaran: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Pembayaran berhasil diverifikasi", nil)
}