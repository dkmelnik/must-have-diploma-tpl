package withdrawals

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Save(ctx context.Context, order *models.Withdrawal) error {
	query := `
		INSERT INTO withdrawals (user_id, order_number, amount)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		order.UserID,
		order.OrderNumber,
		order.Amount,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) IsEntryByOrderNumber(ctx context.Context, orderNumber string) (bool, error) {
	var id string
	query := `
		SELECT id FROM withdrawals WHERE order_number = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, orderNumber).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) FindSumOfAmounts(ctx context.Context, userID models.ModelID) (float64, error) {
	var total float64

	query := `
		SELECT COALESCE(SUM(amount), 0) 
		FROM withdrawals
		where user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&total)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.ErrNotFound
		}
		return 0, err
	}

	return total, nil
}

func (r *Repository) Find(ctx context.Context, userID models.ModelID) ([]*models.Withdrawal, error) {
	query := `
		SELECT id, user_id, order_number, amount, created_at
		FROM withdrawals
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Withdrawal
	for rows.Next() {
		var order models.Withdrawal
		if err := rows.Scan(&order.ID, &order.UserID, &order.OrderNumber, &order.Amount, &order.CreatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
