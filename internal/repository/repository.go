package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-notifier/internal/model"
)

var (
	ErrLockTaken               = errors.New("lock is already taken")
	ErrEmailIsNotUnique        = errors.New("email is not unique")
	ErrUserNotFound            = errors.New("user not found")
	ErrSubscriptionIsNotUnique = errors.New("subscription is not unique")
	ErrSubscriptionNotFound    = errors.New("subscription not found")
	ErrQueryResultUnknown      = errors.New("query result unknown")
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, name, surname, passwordHash string, birthdayDate time.Time) (string, error)
	FindByEmail(ctx context.Context, email string) (model.User, error)
	FindAllUsers(ctx context.Context, userId string) ([]model.User, error)
	FindUsersSubscribedTo(ctx context.Context, userId string) ([]model.User, error)
}

type SubscriptionRepository interface {
	CreateSubscription(ctx context.Context, userId, subscriberId string, notifyBeforeDays int) error
	DeleteSubscription(ctx context.Context, userId, subscriberId string) error
}

type NotificationRepository interface {
	GetLock(ctx context.Context) (*sqlx.Tx, error)
	FindUsersToNotify(ctx context.Context, tx *sqlx.Tx) ([]model.Notification, error)
}

type Repository struct {
	User         UserRepository
	Subscription SubscriptionRepository
	Notification NotificationRepository
}
