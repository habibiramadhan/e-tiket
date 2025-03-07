//internal/delivery/http/handler/user_handler.go

package handler

import (
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/usecase"
	"ticket-system/pkg/utils"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(userUsecase usecase.UserUsecase) *UserHandler {
	return &UserHandler{
		userUsecase: userUsecase,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req usecase.RegisterRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	// Validasi input
	var validationErrors []utils.ErrorDetail
	
	if req.Username == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "username",
			Message: "Username tidak boleh kosong",
		})
	}
	
	if req.Email == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "email",
			Message: "Email tidak boleh kosong",
		})
	}
	
	if req.Password == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "password",
			Message: "Password tidak boleh kosong",
		})
	}
	
	if req.Password != req.RetypePassword {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "retype_password",
			Message: "Password dan konfirmasi password tidak cocok",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	userID, err := h.userUsecase.Register(c.Context(), req)
	if err != nil {
		// Deteksi jenis error dan kirim kode error yang sesuai
		switch err.Error() {
		case "username sudah digunakan":
			return utils.ErrorResponse(c, utils.ErrorCodeDuplicateUsername, "Username sudah digunakan", fiber.StatusBadRequest)
		case "email sudah digunakan":
			return utils.ErrorResponse(c, utils.ErrorCodeDuplicateEmail, "Email sudah digunakan", fiber.StatusBadRequest)
		case "format email tidak valid":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidFormat, "Format email tidak valid", fiber.StatusBadRequest)
		case "username harus minimal 4 karakter":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidFormat, "Username harus minimal 4 karakter", fiber.StatusBadRequest)
		case "password harus minimal 6 karakter":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidFormat, "Password harus minimal 6 karakter", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal melakukan registrasi: "+err.Error())
		}
	}
	
	return utils.CreatedResponse(c, "Registrasi berhasil. Silakan cek email Anda untuk verifikasi.", fiber.Map{
		"user_id": userID,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req usecase.LoginRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	if req.Username == "" {
		return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "Username atau email tidak boleh kosong", fiber.StatusBadRequest)
	}
	
	if req.Password == "" {
		return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "Password tidak boleh kosong", fiber.StatusBadRequest)
	}
	
	resp, err := h.userUsecase.Login(c.Context(), req)
	if err != nil {
		switch err.Error() {
		case "username atau password salah":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidCredentials, "Username atau password salah", fiber.StatusUnauthorized)
		case "email belum diverifikasi, silakan periksa email Anda":
			return utils.ErrorResponse(c, utils.ErrorCodeEmailNotVerified, "Email belum diverifikasi, silakan periksa email Anda", fiber.StatusUnauthorized)
		default:
			return utils.ServerError(c, "Gagal melakukan login: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Login berhasil", resp)
}

func (h *UserHandler) CreateProfile(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	var profile struct {
		Name        string `json:"name"`
		Gender      string `json:"gender"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phone_number"`
	}
	
	if err := c.BodyParser(&profile); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	// Validasi input
	var validationErrors []utils.ErrorDetail
	
	if profile.Name == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "name",
			Message: "Nama tidak boleh kosong",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	// Panggil usecase untuk create profile
	profileID, err := h.userUsecase.CreateProfile(c.Context(), userID, profile.Name, profile.Gender, profile.Address, profile.PhoneNumber)
	if err != nil {
		if err.Error() == "user sudah memiliki profil" {
			return utils.ErrorResponse(c, utils.ErrorCodeResourceAlreadyExist, "Profil sudah dibuat sebelumnya", fiber.StatusBadRequest)
		}
		return utils.ServerError(c, "Gagal membuat profil: "+err.Error())
	}
	
	return utils.SuccessResponse(c, "Profil berhasil dibuat", fiber.Map{
		"profile_id": profileID,
	})
}

func (h *UserHandler) VerifyEmail(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "Parameter token diperlukan", fiber.StatusBadRequest)
	}
	
	err := h.userUsecase.VerifyEmail(c.Context(), token)
	if err != nil {
		switch err.Error() {
		case "token verifikasi tidak valid":
			return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token verifikasi tidak valid", fiber.StatusBadRequest)
		case "token verifikasi sudah kedaluwarsa":
			return utils.ErrorResponse(c, utils.ErrorCodeVerificationExpired, "Token verifikasi sudah kedaluwarsa", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal verifikasi email: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Email berhasil diverifikasi", nil)
}

func (h *UserHandler) ResendVerificationEmail(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
	}
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	if req.Email == "" {
		return utils.ErrorResponse(c, utils.ErrorCodeMissingRequiredField, "Email tidak boleh kosong", fiber.StatusBadRequest)
	}
	
	err := h.userUsecase.ResendVerificationEmail(c.Context(), req.Email)
	if err != nil {
		switch err.Error() {
		case "email tidak terdaftar":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Email tidak terdaftar", fiber.StatusBadRequest)
		case "email sudah diverifikasi":
			return utils.ErrorResponse(c, utils.ErrorCodeEmailAlreadyVerified, "Email sudah diverifikasi", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal mengirim ulang email verifikasi: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Email verifikasi berhasil dikirim ulang", nil)
}