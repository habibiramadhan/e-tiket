//internal/domain/entity/event.go

package entity

import "time"

type Event struct {
	ID          int       `json:"id"`
	OwnerID     int       `json:"owner_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	EventDate   time.Time `json:"event_date"`
	MaxCapacity int       `json:"max_capacity"`
	TicketsSold int       `json:"tickets_sold"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}