//internal/domain/entity/transaction.go

package entity

import "time"

type Transaction struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	EventID         int       `json:"event_id"`
	TransactionCode string    `json:"transaction_code"`
	Quantity        int       `json:"quantity"`
	TotalAmount     float64   `json:"total_amount"`
	Status          string    `json:"status"` 
	PaymentMethod   string    `json:"payment_method"`
	PaymentDetail   string    `json:"payment_detail"`
	PaymentProof    string    `json:"payment_proof"`
	VerifiedAt      time.Time `json:"verified_at,omitempty"`
	VerifiedBy      int       `json:"verified_by,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}