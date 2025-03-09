//internal/repository/postgres/transaction_repository.go

package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"
	"ticket-system/internal/domain/entity"
)

type transactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *transactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) (int, error) {
	query := `
		INSERT INTO transactions (
			user_id, event_id, transaction_code, quantity, total_amount, 
			status, payment_method, payment_detail, payment_proof,
			created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id
	`

	var id int
	err := r.db.QueryRowContext(
		ctx,
		query,
		transaction.UserID,
		transaction.EventID,
		transaction.TransactionCode,
		transaction.Quantity,
		transaction.TotalAmount,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.PaymentDetail,
		transaction.PaymentProof,
		transaction.CreatedAt,
		transaction.UpdatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *transactionRepository) FindByID(ctx context.Context, id int) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, event_id, transaction_code, quantity, 
			total_amount, status, payment_method, payment_detail, payment_proof,
			verified_at, verified_by, created_at, updated_at
		FROM transactions
		WHERE id = $1
	`

	var transaction entity.Transaction
	var verifiedAt sql.NullTime
	var verifiedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.EventID,
		&transaction.TransactionCode,
		&transaction.Quantity,
		&transaction.TotalAmount,
		&transaction.Status,
		&transaction.PaymentMethod,
		&transaction.PaymentDetail,
		&transaction.PaymentProof,
		&verifiedAt,
		&verifiedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if verifiedAt.Valid {
		transaction.VerifiedAt = verifiedAt.Time
	}
	if verifiedBy.Valid {
		transaction.VerifiedBy = int(verifiedBy.Int64)
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByCode(ctx context.Context, code string) (*entity.Transaction, error) {
	query := `
		SELECT id, user_id, event_id, transaction_code, quantity, 
			total_amount, status, payment_method, payment_detail, payment_proof,
			verified_at, verified_by, created_at, updated_at
		FROM transactions
		WHERE transaction_code = $1
	`

	var transaction entity.Transaction
	var verifiedAt sql.NullTime
	var verifiedBy sql.NullInt64

	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&transaction.ID,
		&transaction.UserID,
		&transaction.EventID,
		&transaction.TransactionCode,
		&transaction.Quantity,
		&transaction.TotalAmount,
		&transaction.Status,
		&transaction.PaymentMethod,
		&transaction.PaymentDetail,
		&transaction.PaymentProof,
		&verifiedAt,
		&verifiedBy,
		&transaction.CreatedAt,
		&transaction.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	if verifiedAt.Valid {
		transaction.VerifiedAt = verifiedAt.Time
	}
	if verifiedBy.Valid {
		transaction.VerifiedBy = int(verifiedBy.Int64)
	}

	return &transaction, nil
}

func (r *transactionRepository) FindByUserID(ctx context.Context, userID, offset, limit int) ([]entity.Transaction, error) {
	query := `
		SELECT id, user_id, event_id, transaction_code, quantity, 
			total_amount, status, payment_method, payment_detail, payment_proof,
			verified_at, verified_by, created_at, updated_at
		FROM transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.QueryContext(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []entity.Transaction
	for rows.Next() {
		var transaction entity.Transaction
		var verifiedAt sql.NullTime
		var verifiedBy sql.NullInt64

		err := rows.Scan(
			&transaction.ID,
			&transaction.UserID,
			&transaction.EventID,
			&transaction.TransactionCode,
			&transaction.Quantity,
			&transaction.TotalAmount,
			&transaction.Status,
			&transaction.PaymentMethod,
			&transaction.PaymentDetail,
			&transaction.PaymentProof,
			&verifiedAt,
			&verifiedBy,
			&transaction.CreatedAt,
			&transaction.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if verifiedAt.Valid {
			transaction.VerifiedAt = verifiedAt.Time
		}
		if verifiedBy.Valid {
			transaction.VerifiedBy = int(verifiedBy.Int64)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func (r *transactionRepository) CountByUserID(ctx context.Context, userID int) (int, error) {
	query := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`

	var count int
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *transactionRepository) Update(ctx context.Context, transaction *entity.Transaction) error {
	query := `
		UPDATE transactions
		SET user_id = $1, event_id = $2, transaction_code = $3, quantity = $4, 
			total_amount = $5, status = $6, payment_method = $7, payment_detail = $8, 
			payment_proof = $9, verified_at = $10, verified_by = $11, updated_at = $12
		WHERE id = $13
	`

	verifiedAt := sql.NullTime{}
	if !transaction.VerifiedAt.IsZero() {
		verifiedAt.Time = transaction.VerifiedAt
		verifiedAt.Valid = true
	}

	verifiedBy := sql.NullInt64{}
	if transaction.VerifiedBy != 0 {
		verifiedBy.Int64 = int64(transaction.VerifiedBy)
		verifiedBy.Valid = true
	}

	_, err := r.db.ExecContext(
		ctx,
		query,
		transaction.UserID,
		transaction.EventID,
		transaction.TransactionCode,
		transaction.Quantity,
		transaction.TotalAmount,
		transaction.Status,
		transaction.PaymentMethod,
		transaction.PaymentDetail,
		transaction.PaymentProof,
		verifiedAt,
		verifiedBy,
		time.Now(),
		transaction.ID,
	)

	return err
}

func (r *transactionRepository) UpdateStatus(ctx context.Context, id int, status string) error {
	query := `UPDATE transactions SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *transactionRepository) UpdatePaymentProof(ctx context.Context, id int, proofURL string) error {
	query := `UPDATE transactions SET payment_proof = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, proofURL, time.Now(), id)
	return err
}

func (r *transactionRepository) VerifyPayment(ctx context.Context, id, verifierID int) error {
	query := `
		UPDATE transactions 
		SET status = 'success', verified_at = $1, verified_by = $2, updated_at = $3
		WHERE id = $4
	`
	_, err := r.db.ExecContext(ctx, query, time.Now(), verifierID, time.Now(), id)
	return err
}