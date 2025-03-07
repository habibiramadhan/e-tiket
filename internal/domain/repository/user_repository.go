//internal/domain/repository/user_repository.go

package repository

import (
	"context"
	"ticket-system/internal/domain/entity"
)

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (int, error)
	FindByID(ctx context.Context, id int) (*entity.User, error)
	FindByUsername(ctx context.Context, username string) (*entity.User, error)
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, id int) error
	UpdateVerificationStatus(ctx context.Context, userID int, isVerified bool) error
}