//internal/domain/entity/user_profile.go

package entity

import (
	"time"
)

type UserProfile struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Gender      string    `json:"gender"`
	Address     string    `json:"address"`
	PhoneNumber string    `json:"phone_number"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}