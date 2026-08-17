package models

import (
	"time"

	"github.com/shopspring/decimal"
)

type Subscription struct {
	ID        int64           `json:"id"`
	UserID    int64           `json:"user_id"`
	Service   string          `json:"service"`
	Price     decimal.Decimal `json:"price"`
	StartDate time.Time       `json:"start_date"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}

type SubscriptionInput struct {
	UserID    int64           `json:"user_id"`
	Service   string          `json:"service"`
	Price     decimal.Decimal `json:"price"`
	StartDate time.Time       `json:"start_date"`
}

type SubscriptionUpdateInput struct {
	Service   string          `json:"service"`
	Price     decimal.Decimal `json:"price"`
	StartDate time.Time       `json:"start_date"`
}
