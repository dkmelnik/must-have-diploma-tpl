package users

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

func (r *Repository) Save(ctx context.Context, m *models.User) (models.ModelID, error) {
	var id models.ModelID
	query := `
		INSERT INTO users (login, password) 
		VALUES ($1, $2) 
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query, m.Login, m.Password).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (r *Repository) IsEntryByLogin(ctx context.Context, login string) (bool, error) {
	var id string
	query := `
		SELECT id FROM users WHERE login = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, login).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func (r *Repository) FindOneByLogin(ctx context.Context, login string) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, login, password, created_at FROM users WHERE login = $1 LIMIT 1
	`
	err := r.db.QueryRowContext(ctx, query, login).Scan(&user.ID, &user.Login, &user.Password, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}
	return &user, nil
}
