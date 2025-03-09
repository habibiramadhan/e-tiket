//internal/domain/repository/email_verification_repository.go

package repository

import (
	"context"
	"ticket-system/internal/domain/entity"
)

type EmailVerificationRepository interface {
	Create(ctx context.Context, verification *entity.EmailVerification) (int, error)
	FindByToken(ctx context.Context, token string) (*entity.EmailVerification, error)
	Delete(ctx context.Context, id int) error
	DeleteByUserID(ctx context.Context, userID int) error
}