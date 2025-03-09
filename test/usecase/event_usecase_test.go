//test/usecase/event_usecase_test.go

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

func TestCreateEvent(t *testing.T) {
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	eventUsecase := usecase.NewEventUsecase(mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		userID := 1
		user := &entity.User{
			ID:       userID,
			Username: "organizer1",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		req := usecase.CreateEventRequest{
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			Price:       250000,
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("Create", ctx, mock.AnythingOfType("*entity.Event")).Return(1, nil).Once()
		
		eventID, err := eventUsecase.CreateEvent(ctx, userID, req)
		
		assert.NoError(t, err)
		assert.Equal(t, 1, eventID)
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("User Not Found", func(t *testing.T) {
		userID := 2
		
		req := usecase.CreateEventRequest{
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			Price:       250000,
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(nil, nil).Once()
		
		eventID, err := eventUsecase.CreateEvent(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, 0, eventID)
		assert.Equal(t, "pengguna tidak ditemukan", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Not An Organizer", func(t *testing.T) {
		userID := 3
		user := &entity.User{
			ID:       userID,
			Username: "regular_user",
			Email:    "user@example.com",
			Role:     "user",
		}
		
		req := usecase.CreateEventRequest{
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			Price:       250000,
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		
		eventID, err := eventUsecase.CreateEvent(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, 0, eventID)
		assert.Equal(t, "hanya organizer yang dapat membuat event", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Past Event Date", func(t *testing.T) {
		userID := 1
		user := &entity.User{
			ID:       userID,
			Username: "organizer1",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		req := usecase.CreateEventRequest{
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(-24 * time.Hour),
			MaxCapacity: 1000,
			Price:       250000,
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		
		eventID, err := eventUsecase.CreateEvent(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, 0, eventID)
		assert.Equal(t, "tanggal event tidak boleh di masa lalu", err.Error())
		mockUserRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error", func(t *testing.T) {
		userID := 1
		user := &entity.User{
			ID:       userID,
			Username: "organizer1",
			Email:    "organizer@example.com",
			Role:     "organizer",
		}
		
		req := usecase.CreateEventRequest{
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			Price:       250000,
		}
		
		mockUserRepo.On("FindByID", ctx, userID).Return(user, nil).Once()
		mockEventRepo.On("Create", ctx, mock.AnythingOfType("*entity.Event")).Return(0, errors.New("database error")).Once()
		
		eventID, err := eventUsecase.CreateEvent(ctx, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, 0, eventID)
		assert.Equal(t, "database error", err.Error())
		mockUserRepo.AssertExpectations(t)
		mockEventRepo.AssertExpectations(t)
	})
}

func TestGetEventByID(t *testing.T) {
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	eventUsecase := usecase.NewEventUsecase(mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		eventID := 1
		event := &entity.Event{
			ID:          eventID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		result, err := eventUsecase.GetEventByID(ctx, eventID)
		
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, event.ID, result.ID)
		assert.Equal(t, event.Title, result.Title)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found", func(t *testing.T) {
		eventID := 999
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, nil).Once()
		
		result, err := eventUsecase.GetEventByID(ctx, eventID)
		
		assert.NoError(t, err)
		assert.Nil(t, result)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error", func(t *testing.T) {
		eventID := 1
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, errors.New("database error")).Once()
		
		result, err := eventUsecase.GetEventByID(ctx, eventID)
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database error", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
}

func TestUpdateEvent(t *testing.T) {
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	eventUsecase := usecase.NewEventUsecase(mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		eventID := 1
		userID := 1
		
		existingEvent := &entity.Event{
			ID:          eventID,
			OwnerID:     userID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.UpdateEventRequest{
			Title:       "Konser Musik Rock (Update)",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama Jakarta",
			EventDate:   time.Now().Add(48 * time.Hour),
			MaxCapacity: 1200,
			Price:       300000,
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(existingEvent, nil).Once()
		mockEventRepo.On("Update", ctx, mock.AnythingOfType("*entity.Event")).Return(nil).Once()
		
		err := eventUsecase.UpdateEvent(ctx, eventID, userID, req)
		
		assert.NoError(t, err)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found", func(t *testing.T) {
		eventID := 999
		userID := 1
		
		req := usecase.UpdateEventRequest{
			Title:       "Konser Musik Rock (Update)",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama Jakarta",
			EventDate:   time.Now().Add(48 * time.Hour),
			MaxCapacity: 1200,
			Price:       300000,
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, nil).Once()
		
		err := eventUsecase.UpdateEvent(ctx, eventID, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "event tidak ditemukan", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Not Event Owner", func(t *testing.T) {
		eventID := 1
		ownerID := 1
		userID := 2
		
		existingEvent := &entity.Event{
			ID:          eventID,
			OwnerID:     ownerID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.UpdateEventRequest{
			Title:       "Konser Musik Rock (Update)",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama Jakarta",
			EventDate:   time.Now().Add(48 * time.Hour),
			MaxCapacity: 1200,
			Price:       300000,
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(existingEvent, nil).Once()
		
		err := eventUsecase.UpdateEvent(ctx, eventID, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "anda tidak memiliki izin untuk mengubah event ini", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Capacity Below Sold Tickets", func(t *testing.T) {
		eventID := 1
		userID := 1
		
		existingEvent := &entity.Event{
			ID:          eventID,
			OwnerID:     userID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.UpdateEventRequest{
			Title:       "Konser Musik Rock (Update)",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama Jakarta",
			EventDate:   time.Now().Add(48 * time.Hour),
			MaxCapacity: 400,
			Price:       300000,
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(existingEvent, nil).Once()
		
		err := eventUsecase.UpdateEvent(ctx, eventID, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "kapasitas tidak boleh lebih kecil dari jumlah tiket yang sudah terjual", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Invalid Status", func(t *testing.T) {
		eventID := 1
		userID := 1
		
		existingEvent := &entity.Event{
			ID:          eventID,
			OwnerID:     userID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		req := usecase.UpdateEventRequest{
			Title:       "Konser Musik Rock (Update)",
			Description: "Konser musik rock dengan berbagai band terkenal",
			Location:    "Stadion Utama Jakarta",
			EventDate:   time.Now().Add(48 * time.Hour),
			MaxCapacity: 1200,
			Price:       300000,
			Status:      "invalid_status",
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(existingEvent, nil).Once()
		
		err := eventUsecase.UpdateEvent(ctx, eventID, userID, req)
		
		assert.Error(t, err)
		assert.Equal(t, "status tidak valid", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
}

func TestGetEventSales(t *testing.T) {
	mockEventRepo := new(mocks.MockEventRepository)
	mockUserRepo := new(mocks.MockUserRepository)
	
	eventUsecase := usecase.NewEventUsecase(mockEventRepo, mockUserRepo)
	ctx := context.Background()
	
	t.Run("Success", func(t *testing.T) {
		eventID := 1
		userID := 1
		
		event := &entity.Event{
			ID:          eventID,
			OwnerID:     userID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		sales, err := eventUsecase.GetEventSales(ctx, eventID, userID)
		
		assert.NoError(t, err)
		assert.NotNil(t, sales)
		assert.Equal(t, eventID, sales.EventID)
		assert.Equal(t, "Konser Musik Rock", sales.Title)
		assert.Equal(t, 1000, sales.MaxCapacity)
		assert.Equal(t, 500, sales.TicketsSold)
		assert.Equal(t, 500, sales.AvailableTickets)
		assert.Equal(t, 250000.0, sales.Price)
		assert.Equal(t, 250000.0*500, sales.TotalSales)
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Event Not Found", func(t *testing.T) {
		eventID := 999
		userID := 1
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, nil).Once()
		
		sales, err := eventUsecase.GetEventSales(ctx, eventID, userID)
		
		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.Equal(t, "event tidak ditemukan", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Not Event Owner", func(t *testing.T) {
		eventID := 1
		ownerID := 1
		userID := 2
		
		event := &entity.Event{
			ID:          eventID,
			OwnerID:     ownerID,
			Title:       "Konser Musik Rock",
			Description: "Konser musik rock",
			Location:    "Stadion Utama",
			EventDate:   time.Now().Add(24 * time.Hour),
			MaxCapacity: 1000,
			TicketsSold: 500,
			Price:       250000,
			Status:      "active",
		}
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(event, nil).Once()
		
		sales, err := eventUsecase.GetEventSales(ctx, eventID, userID)
		
		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.Equal(t, "anda tidak memiliki izin untuk melihat data penjualan event ini", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
	
	t.Run("Database Error", func(t *testing.T) {
		eventID := 1
		userID := 1
		
		mockEventRepo.On("FindByID", ctx, eventID).Return(nil, errors.New("database error")).Once()
		
		sales, err := eventUsecase.GetEventSales(ctx, eventID, userID)
		
		assert.Error(t, err)
		assert.Nil(t, sales)
		assert.Equal(t, "database error", err.Error())
		mockEventRepo.AssertExpectations(t)
	})
}