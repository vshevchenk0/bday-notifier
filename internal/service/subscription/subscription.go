package subscription

import (
	"context"
	"errors"
	"log/slog"

	"github.com/vshevchenk0/bday-greeter/internal/repository"
	"github.com/vshevchenk0/bday-greeter/internal/service"
)

type subscriptionService struct {
	subscriptionRepository repository.SubscriptionRepository
	logger                 *slog.Logger
}

func NewSubscriptionService(
	subscriptionRepository repository.SubscriptionRepository,
	logger *slog.Logger,
) *subscriptionService {
	return &subscriptionService{
		subscriptionRepository: subscriptionRepository,
		logger:                 logger,
	}
}

func (s *subscriptionService) CreateSubscription(
	ctx context.Context, userId, subscriberId string, notifyBeforeDays int,
) error {
	err := s.subscriptionRepository.CreateSubscription(ctx, userId, subscriberId, notifyBeforeDays)
	if errors.Is(err, repository.ErrUserNotFound) {
		return service.ErrUserNotFound
	}
	if errors.Is(err, repository.ErrSubscriptionIsNotUnique) {
		return service.ErrDuplicateSubscription
	}
	if err != nil {
		return errors.New("failed to create subscription")
	}
	return nil
}

func (s *subscriptionService) DeleteSubscription(ctx context.Context, userId, subscriberId string) error {
	err := s.subscriptionRepository.DeleteSubscription(ctx, userId, subscriberId)
	if errors.Is(err, repository.ErrSubscriptionNotFound) {
		return service.ErrSubscriptionNotFound
	}
	if errors.Is(err, repository.ErrQueryResultUnknown) {
		return service.ErrOperationResultUnknown
	}
	if err != nil {
		return errors.New("failed to delete subscription")
	}
	return nil
}
