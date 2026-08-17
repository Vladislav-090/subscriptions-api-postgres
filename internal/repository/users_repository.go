package repository

import (
	"context"
	"database/sql"
	"subscriptions-api-postgres/internal/models"
)

type UsersRepository struct {
	db *sql.DB
}

func NewUsersRepository(db *sql.DB) *UsersRepository {
	return &UsersRepository{
		db: db,
	}
}

func (r *UsersRepository) CreateUser(
	ctx context.Context,
	email string,
	passwordHash string,
	role models.Role,
) (models.User, error) {
	var user models.User

	query := `INSERT INTO users (email, password_hash, role)
VALUES ($1, $2, $3)
RETURNING id, email, password_hash, role, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
		passwordHash,
		role,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *UsersRepository) GetUserByEmail(
	ctx context.Context,
	email string,
) (models.User, error) {
	query := `SELECT id, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user models.User
	err := r.db.QueryRowContext(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
