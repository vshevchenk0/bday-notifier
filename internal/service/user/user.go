package user

import (
	"context"
	"errors"
	"log/slog"

	"github.com/vshevchenk0/bday-notifier/internal/model"
	"github.com/vshevchenk0/bday-notifier/internal/repository"
)

type userService struct {
	userRepository repository.UserRepository
	logger         *slog.Logger
}

func NewUserService(userRepository repository.UserRepository, logger *slog.Logger) *userService {
	return &userService{
		userRepository: userRepository,
		logger:         logger,
	}
}

func (s *userService) FindAllUsers(ctx context.Context, userId string) ([]model.User, error) {
	users, err := s.userRepository.FindAllUsers(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to find users")
	}
	return users, nil
}

func (s *userService) FindUsersSubscribedTo(ctx context.Context, userId string) ([]model.User, error) {
	users, err := s.userRepository.FindUsersSubscribedTo(ctx, userId)
	if err != nil {
		return nil, errors.New("failed to find users subscribed to")
	}
	return users, nil
}
