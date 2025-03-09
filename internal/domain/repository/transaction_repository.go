//internal/domain/repository/transaction_repository.go

package repository

import (
	"context"
	"ticket-system/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) (int, error)
	FindByID(ctx context.Context, id int) (*entity.Transaction, error)
	FindByCode(ctx context.Context, code string) (*entity.Transaction, error)
	FindByUserID(ctx context.Context, userID, offset, limit int) ([]entity.Transaction, error)
	CountByUserID(ctx context.Context, userID int) (int, error)
	Update(ctx context.Context, transaction *entity.Transaction) error
	UpdateStatus(ctx context.Context, id int, status string) error
	UpdatePaymentProof(ctx context.Context, id int, proofURL string) error
	VerifyPayment(ctx context.Context, id, verifierID int) error
}