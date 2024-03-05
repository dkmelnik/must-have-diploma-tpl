package orders

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/dkmelnik/go-musthave-diploma/internal/apperrors"
	"github.com/dkmelnik/go-musthave-diploma/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db}
}

func (r *Repository) Save(ctx context.Context, order *models.Order) error {
	query := `
		INSERT INTO orders (user_id, number, status)
		VALUES ($1, $2, $3)
		RETURNING id
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		order.UserID,
		order.Number,
		order.Status,
	)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateByNumber(ctx context.Context, order *models.Order) error {
	query := `
		UPDATE orders
		SET status = $1, accrual = $2
		WHERE number = $3
	`
	_, err := r.db.ExecContext(ctx, query, order.Status, order.Accrual, order.Number)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) FindOneByNumber(ctx context.Context, number string) (*models.Order, error) {
	query := `
		SELECT id, user_id, number, status, created_at, updated_at
		FROM orders
		WHERE number = $1
		LIMIT 1
	`

	var order models.Order
	err := r.db.QueryRowContext(ctx, query, number).Scan(
		&order.ID,
		&order.UserID,
		&order.Number,
		&order.Status,
		&order.CreatedAt,
		&order.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return &order, nil
}

func (r *Repository) FindByUserID(ctx context.Context, userID models.ModelID) ([]*models.Order, error) {
	query := `
			SELECT id, user_id, number, accrual, status, created_at, updated_at 
			FROM orders
			WHERE user_id = $1
			ORDER BY created_at 
		`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.Accrual, &order.Status, &order.CreatedAt, &order.UpdatedAt); err != nil {
			return nil, err
		}
		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *Repository) FindSumOfAccruals(ctx context.Context, userID models.ModelID) (float64, error) {
	var totalAccrual float64

	query := `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		where user_id = $1
	`

	err := r.db.QueryRowContext(ctx, query, userID).Scan(&totalAccrual)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, apperrors.ErrNotFound
		}
		return 0, err
	}

	return totalAccrual, nil
}
