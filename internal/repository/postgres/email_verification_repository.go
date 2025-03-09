//internal/repository/postgres/email_verification_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"ticket-system/internal/domain/entity"
)

type emailVerificationRepository struct {
	db *sql.DB
}

func NewEmailVerificationRepository(db *sql.DB) *emailVerificationRepository {
	return &emailVerificationRepository{
		db: db,
	}
}

func (r *emailVerificationRepository) Create(ctx context.Context, verification *entity.EmailVerification) (int, error) {
	query := `
		INSERT INTO email_verifications (user_id, token, expired_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		verification.UserID,
		verification.Token,
		verification.ExpiredAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *emailVerificationRepository) FindByToken(ctx context.Context, token string) (*entity.EmailVerification, error) {
	query := `
		SELECT id, user_id, token, expired_at, created_at
		FROM email_verifications
		WHERE token = $1
	`

	var verification entity.EmailVerification
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&verification.ID,
		&verification.UserID,
		&verification.Token,
		&verification.ExpiredAt,
		&verification.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &verification, nil
}

func (r *emailVerificationRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM email_verifications WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *emailVerificationRepository) DeleteByUserID(ctx context.Context, userID int) error {
	query := `DELETE FROM email_verifications WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}