//internal/repository/postgres/event_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	
	"ticket-system/internal/domain/entity"
)

type eventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) *eventRepository {
	return &eventRepository{
		db: db,
	}
}

func (r *eventRepository) Create(ctx context.Context, event *entity.Event) (int, error) {
	query := `
		INSERT INTO events (owner_id, title, description, location, event_date, max_capacity, tickets_sold, price, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`
	
	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		event.OwnerID,
		event.Title,
		event.Description,
		event.Location,
		event.EventDate,
		event.MaxCapacity,
		event.TicketsSold,
		event.Price,
		event.Status,
		event.CreatedAt,
		event.UpdatedAt,
	).Scan(&id)
	
	if err != nil {
		return 0, err
	}
	
	return id, nil
}

func (r *eventRepository) FindByID(ctx context.Context, id int) (*entity.Event, error) {
	query := `
		SELECT id, owner_id, title, description, location, event_date, max_capacity, tickets_sold, price, status, created_at, updated_at
		FROM events
		WHERE id = $1
	`
	
	var event entity.Event
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&event.ID,
		&event.OwnerID,
		&event.Title,
		&event.Description,
		&event.Location,
		&event.EventDate,
		&event.MaxCapacity,
		&event.TicketsSold,
		&event.Price,
		&event.Status,
		&event.CreatedAt,
		&event.UpdatedAt,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	
	return &event, nil
}

func (r *eventRepository) FindAll(ctx context.Context, offset, limit int) ([]entity.Event, error) {
	query := `
		SELECT id, owner_id, title, description, location, event_date, max_capacity, tickets_sold, price, status, created_at, updated_at
		FROM events
		WHERE status = 'active'
		ORDER BY event_date ASC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var events []entity.Event
	for rows.Next() {
		var event entity.Event
		err := rows.Scan(
			&event.ID,
			&event.OwnerID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.EventDate,
			&event.MaxCapacity,
			&event.TicketsSold,
			&event.Price,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	
	return events, nil
}

func (r *eventRepository) CountAll(ctx context.Context) (int, error) {
	query := `SELECT COUNT(*) FROM events WHERE status = 'active'`
	
	var count int
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *eventRepository) FindByOwnerID(ctx context.Context, ownerID, offset, limit int) ([]entity.Event, error) {
	query := `
		SELECT id, owner_id, title, description, location, event_date, max_capacity, tickets_sold, price, status, created_at, updated_at
		FROM events
		WHERE owner_id = $1
		ORDER BY event_date ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.QueryContext(ctx, query, ownerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var events []entity.Event
	for rows.Next() {
		var event entity.Event
		err := rows.Scan(
			&event.ID,
			&event.OwnerID,
			&event.Title,
			&event.Description,
			&event.Location,
			&event.EventDate,
			&event.MaxCapacity,
			&event.TicketsSold,
			&event.Price,
			&event.Status,
			&event.CreatedAt,
			&event.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	
	return events, nil
}

func (r *eventRepository) CountByOwnerID(ctx context.Context, ownerID int) (int, error) {
	query := `SELECT COUNT(*) FROM events WHERE owner_id = $1`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, ownerID).Scan(&count)
	if err != nil {
		return 0, err
	}
	
	return count, nil
}

func (r *eventRepository) Update(ctx context.Context, event *entity.Event) error {
	query := `
		UPDATE events
		SET title = $1, description = $2, location = $3, event_date = $4, max_capacity = $5, tickets_sold = $6, price = $7, status = $8, updated_at = $9
		WHERE id = $10
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		event.Title,
		event.Description,
		event.Location,
		event.EventDate,
		event.MaxCapacity,
		event.TicketsSold,
		event.Price,
		event.Status,
		time.Now(),
		event.ID,
	)
	
	return err
}

func (r *eventRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM events WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *eventRepository) UpdateTicketsSold(ctx context.Context, eventID, quantity int) error {
	query := `
		UPDATE events
		SET tickets_sold = tickets_sold + $1, updated_at = $2
		WHERE id = $3
	`
	
	_, err := r.db.ExecContext(
		ctx,
		query,
		quantity,
		time.Now(),
		eventID,
	)
	
	return err
}