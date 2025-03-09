//internal/delivery/http/handler/event_handler.go

package handler

import (
	"strconv"
	"github.com/gofiber/fiber/v2"
	
	"ticket-system/internal/usecase"
	"ticket-system/pkg/utils"
)

type EventHandler struct {
	eventUsecase usecase.EventUsecase
}

func NewEventHandler(eventUsecase usecase.EventUsecase) *EventHandler {
	return &EventHandler{
		eventUsecase: eventUsecase,
	}
}

func (h *EventHandler) CreateEvent(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}

	if claims.Role != "organizer" {
		return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Anda tidak memiliki izin untuk membuat event", fiber.StatusForbidden)
	}
	
	var req usecase.CreateEventRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	var validationErrors []utils.ErrorDetail
	
	if req.Title == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "title",
			Message: "Judul event tidak boleh kosong",
		})
	}
	
	if req.EventDate.IsZero() {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "event_date",
			Message: "Tanggal event tidak boleh kosong",
		})
	}
	
	if req.MaxCapacity <= 0 {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "max_capacity",
			Message: "Kapasitas maksimal harus lebih dari 0",
		})
	}
	
	if req.Price < 0 {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "price",
			Message: "Harga tiket tidak boleh negatif",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	eventID, err := h.eventUsecase.CreateEvent(c.Context(), userID, req)
	if err != nil {
		switch err.Error() {
		case "pengguna tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeResourceNotFound, "Pengguna tidak ditemukan", fiber.StatusNotFound)
		case "hanya organizer yang dapat membuat event":
			return utils.ErrorResponse(c, utils.ErrorCodeUnauthorized, "Hanya organizer yang dapat membuat event", fiber.StatusForbidden)
		case "tanggal event tidak boleh di masa lalu":
			return utils.ErrorResponse(c, utils.ErrorCodeEventDateInvalid, "Tanggal event tidak boleh di masa lalu", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal membuat event: "+err.Error())
		}
	}
	
	event, err := h.eventUsecase.GetEventByID(c.Context(), eventID)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan detail event: "+err.Error())
	}
	
	return utils.CreatedResponse(c, "Event berhasil dibuat", event)
}

func (h *EventHandler) GetEventList(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}
	
	events, total, err := h.eventUsecase.GetEventList(c.Context(), page, limit)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan daftar event: "+err.Error())
	}
	
	meta := fiber.Map{
		"page":  page,
		"limit": limit,
		"total": total,
	}
	
	return utils.SuccessResponse(c, "Daftar event berhasil diambil", events, meta)
}

func (h *EventHandler) GetEventByID(c *fiber.Ctx) error {
	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID event tidak valid", fiber.StatusBadRequest)
	}
	
	event, err := h.eventUsecase.GetEventByID(c.Context(), eventID)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan detail event: "+err.Error())
	}
	
	if event == nil {
		return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event tidak ditemukan", fiber.StatusNotFound)
	}
	
	return utils.SuccessResponse(c, "Detail event berhasil diambil", event)
}

func (h *EventHandler) UpdateEvent(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID event tidak valid", fiber.StatusBadRequest)
	}
	
	var req usecase.UpdateEventRequest
	
	if err := c.BodyParser(&req); err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Format JSON tidak valid", fiber.StatusBadRequest)
	}
	
	var validationErrors []utils.ErrorDetail
	
	if req.Title == "" {
		validationErrors = append(validationErrors, utils.ErrorDetail{
			Field:   "title",
			Message: "Judul event tidak boleh kosong",
		})
	}
	
	if len(validationErrors) > 0 {
		return utils.ValidationError(c, "Validasi gagal", validationErrors)
	}
	
	err = h.eventUsecase.UpdateEvent(c.Context(), eventID, userID, req)
	if err != nil {
		switch err.Error() {
		case "event tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk mengubah event ini":
			return utils.ErrorResponse(c, utils.ErrorCodeEventOwnership, "Anda tidak memiliki izin untuk mengubah event ini", fiber.StatusForbidden)
		case "kapasitas tidak boleh lebih kecil dari jumlah tiket yang sudah terjual":
			return utils.ErrorResponse(c, utils.ErrorCodeEventCapacityLow, "Kapasitas tidak boleh lebih kecil dari jumlah tiket yang sudah terjual", fiber.StatusBadRequest)
		case "status tidak valid":
			return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "Status event tidak valid", fiber.StatusBadRequest)
		default:
			return utils.ServerError(c, "Gagal mengubah event: "+err.Error())
		}
	}
	
	event, err := h.eventUsecase.GetEventByID(c.Context(), eventID)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan detail event: "+err.Error())
	}
	
	return utils.SuccessResponse(c, "Event berhasil diperbarui", event)
}

func (h *EventHandler) DeleteEvent(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID event tidak valid", fiber.StatusBadRequest)
	}
	
	err = h.eventUsecase.DeleteEvent(c.Context(), eventID, userID)
	if err != nil {
		switch err.Error() {
		case "event tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk menghapus event ini":
			return utils.ErrorResponse(c, utils.ErrorCodeEventOwnership, "Anda tidak memiliki izin untuk menghapus event ini", fiber.StatusForbidden)
		default:
			return utils.ServerError(c, "Gagal menghapus event: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Event berhasil dihapus", nil)
}

func (h *EventHandler) GetEventsByOrganizer(c *fiber.Ctx) error {
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
	
	events, total, err := h.eventUsecase.GetEventsByOrganizer(c.Context(), userID, page, limit)
	if err != nil {
		return utils.ServerError(c, "Gagal mendapatkan daftar event: "+err.Error())
	}
	
	meta := fiber.Map{
		"page":  page,
		"limit": limit,
		"total": total,
	}
	
	return utils.SuccessResponse(c, "Daftar event berhasil diambil", events, meta)
}

func (h *EventHandler) GetEventSales(c *fiber.Ctx) error {
	claims := c.Locals("claims").(*utils.JWTClaim)
	
	userID, err := utils.GetUserIDFromToken(claims)
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeTokenInvalid, "Token tidak valid", fiber.StatusUnauthorized)
	}
	
	eventID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return utils.ErrorResponse(c, utils.ErrorCodeInvalidInput, "ID event tidak valid", fiber.StatusBadRequest)
	}
	
	sales, err := h.eventUsecase.GetEventSales(c.Context(), eventID, userID)
	if err != nil {
		switch err.Error() {
		case "event tidak ditemukan":
			return utils.ErrorResponse(c, utils.ErrorCodeEventNotFound, "Event tidak ditemukan", fiber.StatusNotFound)
		case "anda tidak memiliki izin untuk melihat data penjualan event ini":
			return utils.ErrorResponse(c, utils.ErrorCodeEventOwnership, "Anda tidak memiliki izin untuk melihat data penjualan event ini", fiber.StatusForbidden)
		default:
			return utils.ServerError(c, "Gagal mendapatkan data penjualan event: "+err.Error())
		}
	}
	
	return utils.SuccessResponse(c, "Data penjualan event berhasil diambil", sales)
}