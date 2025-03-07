//internal/repository/postgres/user_profile_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"ticket-system/internal/domain/entity"
)

type userProfileRepository struct {
	db *sql.DB
}

func NewUserProfileRepository(db *sql.DB) *userProfileRepository {
	return &userProfileRepository{
		db: db,
	}
}

func (r *userProfileRepository) Create(ctx context.Context, profile *entity.UserProfile) (int, error) {
	query := `
		INSERT INTO user_profiles (user_id, name, gender, address, phone_number, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		profile.UserID,
		profile.Name,
		profile.Gender,
		profile.Address,
		profile.PhoneNumber,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *userProfileRepository) FindByUserID(ctx context.Context, userID int) (*entity.UserProfile, error) {
	query := `
		SELECT id, user_id, name, gender, address, phone_number, created_at, updated_at
		FROM user_profiles
		WHERE user_id = $1
	`

	var profile entity.UserProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&profile.ID,
		&profile.UserID,
		&profile.Name,
		&profile.Gender,
		&profile.Address,
		&profile.PhoneNumber,
		&profile.CreatedAt,
		&profile.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &profile, nil
}

func (r *userProfileRepository) Update(ctx context.Context, profile *entity.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET name = $1, gender = $2, address = $3, phone_number = $4, updated_at = NOW()
		WHERE id = $5
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		profile.Name,
		profile.Gender,
		profile.Address,
		profile.PhoneNumber,
		profile.ID,
	)

	return err
}

func (r *userProfileRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM user_profiles WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}