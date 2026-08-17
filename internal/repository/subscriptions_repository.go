package repository

import (
	"context"
	"database/sql"
	"subscriptions-api-postgres/internal/models"
)

type SubscriptionsRepository struct {
	db *sql.DB
}

func NewSubscriptionsRepository(db *sql.DB) *SubscriptionsRepository {
	return &SubscriptionsRepository{
		db: db,
	}
}

func (r *SubscriptionsRepository) CreateSubscription(
	ctx context.Context,
	input models.SubscriptionInput,
) (models.Subscription, error) {
	var subscription models.Subscription

	query := `INSERT INTO subscriptions (user_id, service, price, start_date)
VALUES ($1, $2, $3, $4)
RETURNING id, user_id, service, price, start_date, created_at, updated_at`

	err := r.db.QueryRowContext(
		ctx,
		query,
		input.UserID,
		input.Service,
		input.Price,
		input.StartDate,
	).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Service,
		&subscription.Price,
		&subscription.StartDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return models.Subscription{}, err
	}

	return subscription, nil

}

func (r *SubscriptionsRepository) GetSubscriptions(
	ctx context.Context,
	userID int64,
) ([]models.Subscription, error) {

	query := `
		SELECT id, user_id, service, price, start_date, created_at, updated_at
		FROM subscriptions
		WHERE user_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	subscriptions := make([]models.Subscription, 0)

	for rows.Next() {
		var subscription models.Subscription

		err = rows.Scan(
			&subscription.ID,
			&subscription.UserID,
			&subscription.Service,
			&subscription.Price,
			&subscription.StartDate,
			&subscription.CreatedAt,
			&subscription.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		subscriptions = append(subscriptions, subscription)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subscriptions, nil

}

func (r *SubscriptionsRepository) GetSubscriptionByID(
	ctx context.Context,
	id int64,
	userID int64,
) (models.Subscription, error) {
	query := `SELECT id, user_id, service, price, start_date, created_at, updated_at
		FROM subscriptions
		WHERE id = $1 and user_id = $2
		`

	var subscription models.Subscription
	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
		userID,
	).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Service,
		&subscription.Price,
		&subscription.StartDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return models.Subscription{}, err
	}
	return subscription, nil
}

func (r *SubscriptionsRepository) UpdateSubscription(
	ctx context.Context,
	id int64,
	userID int64,
	input models.SubscriptionUpdateInput,
) (models.Subscription, error) {
	query := `UPDATE subscriptions
		SET service = $1, price = $2, start_date = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4 and user_id = $5
		RETURNING id, user_id, service, price, start_date, created_at, updated_at`

	var subscription models.Subscription
	err := r.db.QueryRowContext(
		ctx,
		query,
		input.Service,
		input.Price,
		input.StartDate,
		id,
		userID,
	).Scan(
		&subscription.ID,
		&subscription.UserID,
		&subscription.Service,
		&subscription.Price,
		&subscription.StartDate,
		&subscription.CreatedAt,
		&subscription.UpdatedAt,
	)
	if err != nil {
		return models.Subscription{}, err
	}
	return subscription, nil
}

func (r *SubscriptionsRepository) DeleteSubscription(
	ctx context.Context,
	id int64,
	userID int64,
) error {
	query := `DELETE FROM subscriptions WHERE id = $1 and user_id = $2`

	result, err := r.db.ExecContext(ctx, query, id, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
