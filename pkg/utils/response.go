//pkg/utils/response.go

package utils

import (
	"github.com/gofiber/fiber/v2"
)

const (
	// Success codes
	StatusSuccess = "SUCCESS"

	// Error codes - Authentication & Authorization
	ErrorCodeInvalidCredentials   = "AUTH001" // Username/email atau password tidak valid
	ErrorCodeTokenInvalid         = "AUTH002" // Token tidak valid atau sudah kadaluarsa
	ErrorCodeTokenMissing         = "AUTH003" // Token tidak ada
	ErrorCodeUnauthorized         = "AUTH004" // Tidak memiliki izin
	ErrorCodeEmailNotVerified     = "AUTH005" // Email belum diverifikasi
	ErrorCodeEmailAlreadyVerified = "AUTH006" // Email sudah diverifikasi
	ErrorCodeVerificationExpired  = "AUTH007" // Token verifikasi sudah kadaluarsa

	// Error codes - Validation
	ErrorCodeInvalidInput         = "VAL001" // Input tidak valid secara umum
	ErrorCodeMissingRequiredField = "VAL002" // Field wajib tidak ada
	ErrorCodeInvalidFormat        = "VAL003" // Format tidak valid (email, password, dll)
	ErrorCodePasswordMismatch     = "VAL004" // Password dan konfirmasi tidak cocok
	ErrorCodeDuplicateUsername    = "VAL005" // Username sudah digunakan
	ErrorCodeDuplicateEmail       = "VAL006" // Email sudah digunakan

	// Error codes - Resource
	ErrorCodeResourceNotFound     = "RES001" // Resource tidak ditemukan
	ErrorCodeResourceAlreadyExist = "RES002" // Resource sudah ada
	ErrorCodeResourceLimit        = "RES003" // Melebihi batas resource

	// Error codes - Server
	ErrorCodeServerError          = "SRV001" // Error server internal
	ErrorCodeDatabaseError        = "SRV002" // Error database
	ErrorCodeExternalServiceError = "SRV003" // Error layanan eksternal
	ErrorCodeMailServiceError     = "SRV004" // Error layanan email
)

// APIResponse adalah struktur standar untuk semua respons API
type APIResponse struct {
	Status     bool        `json:"status"`           // true untuk sukses, false untuk gagal
	StatusCode string      `json:"status_code"`      // Kode status: SUCCESS atau kode error
	Message    string      `json:"message"`          // Pesan sukses/error untuk user
	Data       interface{} `json:"data,omitempty"`   // Data yang dikembalikan (hanya jika sukses)
	Errors     interface{} `json:"errors,omitempty"` // Detail error (hanya jika error)
	Meta       interface{} `json:"meta,omitempty"`   // Metadata seperti pagination
}

// ErrorDetail berisi detail tentang error
type ErrorDetail struct {
	Field   string `json:"field,omitempty"`   // Field yang menyebabkan error
	Message string `json:"message,omitempty"` // Pesan error detail
}

// SuccessResponse mengirimkan respons sukses standar
func SuccessResponse(c *fiber.Ctx, message string, data interface{}, meta ...interface{}) error {
	var metaData interface{}
	if len(meta) > 0 {
		metaData = meta[0]
	}

	return c.Status(fiber.StatusOK).JSON(APIResponse{
		Status:     true,
		StatusCode: StatusSuccess,
		Message:    message,
		Data:       data,
		Meta:       metaData,
	})
}

// CreatedResponse mengirimkan respons sukses 201 Created
func CreatedResponse(c *fiber.Ctx, message string, data interface{}) error {
	return c.Status(fiber.StatusCreated).JSON(APIResponse{
		Status:     true,
		StatusCode: StatusSuccess,
		Message:    message,
		Data:       data,
	})
}

// ErrorResponse mengirimkan respons error standar
func ErrorResponse(c *fiber.Ctx, statusCode string, message string, httpStatus int, errors ...interface{}) error {
	var errorDetails interface{}
	if len(errors) > 0 {
		errorDetails = errors[0]
	}

	return c.Status(httpStatus).JSON(APIResponse{
		Status:     false,
		StatusCode: statusCode,
		Message:    message,
		Errors:     errorDetails,
	})
}

// ValidationError mengirimkan respons error validasi
func ValidationError(c *fiber.Ctx, message string, errors []ErrorDetail) error {
	return ErrorResponse(c, ErrorCodeInvalidInput, message, fiber.StatusBadRequest, errors)
}

// UnauthorizedError mengirimkan respons error tidak terautentikasi
func UnauthorizedError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrorCodeUnauthorized, message, fiber.StatusUnauthorized)
}

// NotFoundError mengirimkan respons error tidak ditemukan
func NotFoundError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrorCodeResourceNotFound, message, fiber.StatusNotFound)
}

// ServerError mengirimkan respons error internal server
func ServerError(c *fiber.Ctx, message string) error {
	return ErrorResponse(c, ErrorCodeServerError, message, fiber.StatusInternalServerError)
}