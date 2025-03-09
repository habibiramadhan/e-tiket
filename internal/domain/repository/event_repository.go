//internal/domain/repository/event_repository.go

package repository

import (
	"context"
	"ticket-system/internal/domain/entity"
)

type EventRepository interface {
	Create(ctx context.Context, event *entity.Event) (int, error)
	FindByID(ctx context.Context, id int) (*entity.Event, error)
	FindAll(ctx context.Context, offset, limit int) ([]entity.Event, error)
	CountAll(ctx context.Context) (int, error)
	FindByOwnerID(ctx context.Context, ownerID, offset, limit int) ([]entity.Event, error)
	CountByOwnerID(ctx context.Context, ownerID int) (int, error)
	Update(ctx context.Context, event *entity.Event) error
	Delete(ctx context.Context, id int) error
	UpdateTicketsSold(ctx context.Context, eventID, quantity int) error
}