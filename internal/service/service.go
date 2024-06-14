package service

import (
	"context"
	"errors"
	"time"

	"github.com/vshevchenk0/bday-greeter/internal/model"
)

var (
	ErrDuplicateUser          = errors.New("duplicate user")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrUserNotFound           = errors.New("user not found")
	ErrDuplicateSubscription  = errors.New("duplicate subscription")
	ErrSubscriptionNotFound   = errors.New("subscription not found")
	ErrOperationResultUnknown = errors.New("operation result unknown")
)

type Token struct {
	AccessToken string `json:"access_token"`
}

type AuthService interface {
	SignUp(ctx context.Context, email, password, name, surname string, birthdayDate time.Time) (Token, error)
	SignIn(ctx context.Context, email, password string) (Token, error)
	VerifyToken(ctx context.Context, tokenString string) (string, error)
}

type NotificationService interface {
	NotifyUsers(ctx context.Context) error
}

type SubscriptionService interface {
	CreateSubscription(ctx context.Context, userId, subscriberId string, notifyBeforeDays int) error
	DeleteSubscription(ctx context.Context, userId, subscriberId string) error
}

type UserService interface {
	FindAllUsers(ctx context.Context, userId string) ([]model.User, error)
	FindUsersSubscribedTo(ctx context.Context, userId string) ([]model.User, error)
}

type Service struct {
	Auth         AuthService
	Notification NotificationService
	Subscription SubscriptionService
	User         UserService
}
