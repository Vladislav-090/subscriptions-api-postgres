package repository

import (
	"context"
	"database/sql"
	"time"
)

type RevokedTokensRepository struct {
	db *sql.DB
}

func NewRevokedTokensRepository(db *sql.DB) *RevokedTokensRepository {
	return &RevokedTokensRepository{
		db: db,
	}
}

func (r *RevokedTokensRepository) RevokeToken(
	ctx context.Context,
	jti string,
	expiresAt time.Time,
) error {
	query := `INSERT INTO revoked_tokens (jti, expires_at) VALUES ($1, $2)`

	_, err := r.db.ExecContext(ctx, query, jti, expiresAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *RevokedTokensRepository) IsRevoked(
	ctx context.Context,
	jti string,
) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE jti = $1)`

	var revoked bool
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&revoked)
	if err != nil {
		return false, err
	}

	return revoked, nil
}
