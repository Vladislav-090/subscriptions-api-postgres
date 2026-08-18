package service

import (
	"context"
	"errors"
	"subscriptions-api-postgres/internal/models"

	"github.com/shopspring/decimal"
)

type SubscriptionsRepository interface {
	CreateSubscription(ctx context.Context, input models.SubscriptionInput) (models.Subscription, error)
	GetSubscriptions(ctx context.Context, userID int64, limit int64, offset int64) ([]models.Subscription, error)
	GetSubscriptionByID(ctx context.Context, id int64, userID int64) (models.Subscription, error)
	UpdateSubscription(ctx context.Context, id int64, userID int64, input models.SubscriptionUpdateInput) (models.Subscription, error)
	DeleteSubscription(ctx context.Context, id int64, userID int64) error
}

type SubscriptionService struct {
	subscriptionsRepo SubscriptionsRepository
}

func NewSubscriptionsService(repo SubscriptionsRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionsRepo: repo,
	}
}

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrInvalidPrice  = errors.New("price must be positive")
	ErrEmptyService  = errors.New("service cannot be empty")
	ErrInvalidSubscriptionID = errors.New("invalid subscription id")
)

const (
	defaultSubscriptionsLimit = 20
	maxSubscriptionsLimit     = 100
)

func (s *SubscriptionService) CreateSubscription(
	ctx context.Context,
	input models.SubscriptionInput,
) (models.Subscription, error) {
	if input.UserID <= 0 {
		return models.Subscription{}, ErrInvalidUserID
	}

	if input.Price.LessThanOrEqual(decimal.Zero) {
		return models.Subscription{}, ErrInvalidPrice
	}

	if input.Service == "" {
		return models.Subscription{}, ErrEmptyService
	}

	subscription := models.SubscriptionInput{
		UserID: input.UserID,
		Service: input.Service,
		Price: input.Price,
		StartDate: input.StartDate,
	}
	createdSubscription, err := s.subscriptionsRepo.CreateSubscription(ctx, subscription)
	if err != nil {
		return models.Subscription{}, err
	}

	return createdSubscription, nil
}

func (s *SubscriptionService) GetSubscriptions(ctx context.Context, userID int64, limit int64, offset int64) ([]models.Subscription, error) {
	if userID <= 0 {
		return []models.Subscription{}, ErrInvalidUserID
	}

	if limit <= 0 {
		limit = defaultSubscriptionsLimit
	}
	if limit > maxSubscriptionsLimit {
		limit = maxSubscriptionsLimit
	}
	if offset < 0 {
		offset = 0
	}

	subscriptions, err := s.subscriptionsRepo.GetSubscriptions(ctx, userID, limit, offset)
	if err != nil {
		return []models.Subscription{}, err
	}

	return subscriptions, nil
}

func (s *SubscriptionService) GetSubscriptionByID(
	ctx context.Context,
	id int64,
	userID int64,
	)(models.Subscription, error){
		if id <= 0 {
			return models.Subscription{}, ErrInvalidSubscriptionID
		}

		if userID <= 0 {
			return models.Subscription{}, ErrInvalidUserID
		}

		subscription, err := s.subscriptionsRepo.GetSubscriptionByID(ctx, id, userID)
		if err != nil {
			return models.Subscription{}, err
		}

		return subscription, nil
	}

func (s *SubscriptionService) UpdateSubscription(
	ctx context.Context,
	id int64,
	userID int64,
	input models.SubscriptionUpdateInput,
) (models.Subscription, error) {
	if id <= 0 {
		return models.Subscription{}, ErrInvalidSubscriptionID
	}

	if userID <= 0 {
		return models.Subscription{}, ErrInvalidUserID
	}

	if input.Price.LessThanOrEqual(decimal.Zero) {
		return models.Subscription{}, ErrInvalidPrice
	}

	if input.Service == "" {
		return models.Subscription{}, ErrEmptyService
	}

	updatedSubscription, err := s.subscriptionsRepo.UpdateSubscription(ctx, id, userID, input)
	if err != nil {
		return models.Subscription{}, err
	}

	return updatedSubscription, nil
}

func (s *SubscriptionService) DeleteSubscription(ctx context.Context, id int64, userID int64) error {
	if id <= 0 {
		return ErrInvalidSubscriptionID
	}

	if userID <= 0 {
		return ErrInvalidUserID
	}

	return s.subscriptionsRepo.DeleteSubscription(ctx, id, userID)
}
