package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/vshevchenk0/bday-notifier/internal/repository"
	"github.com/vshevchenk0/bday-notifier/internal/service"
	"github.com/vshevchenk0/bday-notifier/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	userRepository repository.UserRepository
	tokenManager   jwt.Manager
	logger         *slog.Logger
}

func NewAuthService(
	userRepository repository.UserRepository,
	tokenManager jwt.Manager,
	logger *slog.Logger,
) *authService {
	return &authService{
		userRepository: userRepository,
		tokenManager:   tokenManager,
		logger:         logger,
	}
}

func (s *authService) SignUp(
	ctx context.Context, email, password, name, surname string, birthdayDate time.Time,
) (service.Token, error) {
	emptyResponse := service.Token{}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("error during password hashing", slog.String("error", err.Error()))
		return emptyResponse, errors.New("error during sign up")
	}

	userId, err := s.userRepository.CreateUser(ctx, email, name, surname, string(passwordHash), birthdayDate)
	if errors.Is(err, repository.ErrEmailIsNotUnique) {
		return emptyResponse, service.ErrDuplicateUser
	}
	if err != nil {
		s.logger.Error("error during saving user to db", slog.String("error", err.Error()))
		return emptyResponse, errors.New("error during sign up")
	}

	token, err := s.tokenManager.NewToken(userId)
	if err != nil {
		s.logger.Error("error during token creation", slog.String("error", err.Error()))
		return emptyResponse, errors.New("signed up succefully, but failed to automatically authorize. please sign in")
	}

	return service.Token{AccessToken: token}, nil
}

func (s *authService) SignIn(ctx context.Context, email, password string) (service.Token, error) {
	emptyResponse := service.Token{}
	user, err := s.userRepository.FindByEmail(ctx, email)
	// check email or any other field, if it contains default value for that field - user was not found
	if user.Email == "" {
		return emptyResponse, service.ErrUserNotFound
	}
	if err != nil {
		return emptyResponse, errors.New("failed to find user")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return emptyResponse, service.ErrInvalidPassword
	}

	token, err := s.tokenManager.NewToken(user.Id)
	if err != nil {
		s.logger.Error("error during token creation", slog.String("error", err.Error()))
		return emptyResponse, errors.New("failed to authorize")
	}

	return service.Token{AccessToken: token}, nil
}

func (s *authService) VerifyToken(ctx context.Context, tokenString string) (string, error) {
	userId, err := s.tokenManager.ParseToken(tokenString)
	return userId, err
}
