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
	userID int64,
	id int64,
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
