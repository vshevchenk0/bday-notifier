package subscription

import (
	"context"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vshevchenk0/bday-greeter/internal/repository"
)

type subscriptionRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *subscriptionRepository {
	return &subscriptionRepository{
		db: db,
	}
}

func (r *subscriptionRepository) CreateSubscription(ctx context.Context, userId, subscriberId string, notifyBeforeDays int) error {
	query := "INSERT INTO subscriptions (user_id, subscriber_id, notify_before_days) VALUES ($1, $2, $3);"
	_, err := r.db.ExecContext(ctx, query, userId, subscriberId, notifyBeforeDays)
	if err, ok := err.(*pq.Error); ok {
		// check foreign key constraint violation
		if err.Code == "23503" {
			return repository.ErrUserNotFound
		}
		// check unique constraint violation
		if err.Code == "23505" {
			return repository.ErrSubscriptionIsNotUnique
		}
	}
	return err
}

func (r *subscriptionRepository) DeleteSubscription(ctx context.Context, userId, subscriberId string) error {
	query := "DELETE FROM subscriptions WHERE user_id=$1 AND subscriber_id=$2;"
	result, err := r.db.ExecContext(ctx, query, userId, subscriberId)
	if err != nil {
		return errors.New("failed to delete subscription")
	}
	count, err := result.RowsAffected()
	if err != nil {
		return repository.ErrQueryResultUnknown
	}
	if count == 0 {
		return repository.ErrSubscriptionNotFound
	}
	return nil
}
