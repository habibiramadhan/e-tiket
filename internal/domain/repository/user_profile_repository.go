//internal/domain/repository/user_profile_repository.go

package repository

import (
	"context"
	"ticket-system/internal/domain/entity"
)

type UserProfileRepository interface {
	Create(ctx context.Context, profile *entity.UserProfile) (int, error)
	FindByUserID(ctx context.Context, userID int) (*entity.UserProfile, error)
	Update(ctx context.Context, profile *entity.UserProfile) error
	Delete(ctx context.Context, id int) error
}