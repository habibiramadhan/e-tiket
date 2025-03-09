//internal/usecase/event_usecase.go

package usecase

import (
	"context"
	"errors"
	"time"
	
	"ticket-system/internal/domain/entity"
	"ticket-system/internal/domain/repository"
)

type CreateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	MaxCapacity int       `json:"max_capacity"`
	Price       float64   `json:"price"`
}

type UpdateEventRequest struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	MaxCapacity int       `json:"max_capacity"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
}

type EventSalesResponse struct {
	EventID          int     `json:"event_id"`
	Title            string  `json:"title"`
	MaxCapacity      int     `json:"max_capacity"`
	TicketsSold      int     `json:"tickets_sold"`
	AvailableTickets int     `json:"available_tickets"`
	Price            float64 `json:"price"`
	TotalSales       float64 `json:"total_sales"`
	Status           string  `json:"status"`
}

type EventUsecase interface {
	CreateEvent(ctx context.Context, userID int, req CreateEventRequest) (int, error)
	GetEventList(ctx context.Context, page, limit int) ([]entity.Event, int, error)
	GetEventByID(ctx context.Context, id int) (*entity.Event, error)
	UpdateEvent(ctx context.Context, eventID, userID int, req UpdateEventRequest) error
	DeleteEvent(ctx context.Context, eventID, userID int) error
	GetEventsByOrganizer(ctx context.Context, userID, page, limit int) ([]entity.Event, int, error)
	GetEventSales(ctx context.Context, eventID, userID int) (*EventSalesResponse, error)
}

type eventUsecase struct {
	eventRepo repository.EventRepository
	userRepo  repository.UserRepository
}

func NewEventUsecase(eventRepo repository.EventRepository, userRepo repository.UserRepository) EventUsecase {
	return &eventUsecase{
		eventRepo: eventRepo,
		userRepo:  userRepo,
	}
}

func (u *eventUsecase) CreateEvent(ctx context.Context, userID int, req CreateEventRequest) (int, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return 0, err
	}
	
	if user == nil {
		return 0, errors.New("pengguna tidak ditemukan")
	}
	
	if user.Role != "organizer" {
		return 0, errors.New("hanya organizer yang dapat membuat event")
	}
	
	if req.EventDate.Before(time.Now()) {
		return 0, errors.New("tanggal event tidak boleh di masa lalu")
	}
	
	event := &entity.Event{
		OwnerID:     userID,
		Title:       req.Title,
		Description: req.Description,
		Location:    req.Location,
		EventDate:   req.EventDate,
		MaxCapacity: req.MaxCapacity,
		TicketsSold: 0,
		Price:       req.Price,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	
	eventID, err := u.eventRepo.Create(ctx, event)
	if err != nil {
		return 0, err
	}
	
	return eventID, nil
}

func (u *eventUsecase) GetEventList(ctx context.Context, page, limit int) ([]entity.Event, int, error) {
	offset := (page - 1) * limit
	events, err := u.eventRepo.FindAll(ctx, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := u.eventRepo.CountAll(ctx)
	if err != nil {
		return nil, 0, err
	}
	
	return events, total, nil
}

func (u *eventUsecase) GetEventByID(ctx context.Context, id int) (*entity.Event, error) {
	event, err := u.eventRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	return event, nil
}

func (u *eventUsecase) UpdateEvent(ctx context.Context, eventID, userID int, req UpdateEventRequest) error {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	
	if event == nil {
		return errors.New("event tidak ditemukan")
	}
	
	if event.OwnerID != userID {
		return errors.New("anda tidak memiliki izin untuk mengubah event ini")
	}
	
	if req.MaxCapacity < event.TicketsSold {
		return errors.New("kapasitas tidak boleh lebih kecil dari jumlah tiket yang sudah terjual")
	}
	
	if req.Status != "" && req.Status != "active" && req.Status != "cancelled" && req.Status != "completed" {
		return errors.New("status tidak valid")
	}
	
	event.Title = req.Title
	event.Description = req.Description
	event.Location = req.Location
	event.EventDate = req.EventDate
	event.MaxCapacity = req.MaxCapacity
	event.Price = req.Price
	event.UpdatedAt = time.Now()
	
	if req.Status != "" {
		event.Status = req.Status
	}
	
	err = u.eventRepo.Update(ctx, event)
	if err != nil {
		return err
	}
	
	return nil
}

func (u *eventUsecase) DeleteEvent(ctx context.Context, eventID, userID int) error {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return err
	}
	
	if event == nil {
		return errors.New("event tidak ditemukan")
	}
	
	if event.OwnerID != userID {
		return errors.New("anda tidak memiliki izin untuk menghapus event ini")
	}
	
	if event.TicketsSold > 0 {
		event.Status = "cancelled"
		event.UpdatedAt = time.Now()
		return u.eventRepo.Update(ctx, event)
	}
	
	return u.eventRepo.Delete(ctx, eventID)
}

func (u *eventUsecase) GetEventsByOrganizer(ctx context.Context, userID, page, limit int) ([]entity.Event, int, error) {
	offset := (page - 1) * limit
	events, err := u.eventRepo.FindByOwnerID(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	
	total, err := u.eventRepo.CountByOwnerID(ctx, userID)
	if err != nil {
		return nil, 0, err
	}
	
	return events, total, nil
}

func (u *eventUsecase) GetEventSales(ctx context.Context, eventID, userID int) (*EventSalesResponse, error) {
	event, err := u.eventRepo.FindByID(ctx, eventID)
	if err != nil {
		return nil, err
	}
	
	if event == nil {
		return nil, errors.New("event tidak ditemukan")
	}
	
	if event.OwnerID != userID {
		return nil, errors.New("anda tidak memiliki izin untuk melihat data penjualan event ini")
	}
	
	sales := &EventSalesResponse{
		EventID:          event.ID,
		Title:            event.Title,
		MaxCapacity:      event.MaxCapacity,
		TicketsSold:      event.TicketsSold,
		AvailableTickets: event.MaxCapacity - event.TicketsSold,
		Price:            event.Price,
		TotalSales:       float64(event.TicketsSold) * event.Price,
		Status:           event.Status,
	}
	
	return sales, nil
}