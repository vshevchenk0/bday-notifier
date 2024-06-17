package app

import (
	"log/slog"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-notifier/internal/api"
	"github.com/vshevchenk0/bday-notifier/internal/config"
	"github.com/vshevchenk0/bday-notifier/internal/middleware"
	"github.com/vshevchenk0/bday-notifier/internal/repository"
	subscriptionRepository "github.com/vshevchenk0/bday-notifier/internal/repository/subscription"
	userRepository "github.com/vshevchenk0/bday-notifier/internal/repository/user"
	"github.com/vshevchenk0/bday-notifier/internal/server"
	"github.com/vshevchenk0/bday-notifier/internal/service"
	authService "github.com/vshevchenk0/bday-notifier/internal/service/auth"
	subscriptionService "github.com/vshevchenk0/bday-notifier/internal/service/subscription"
	userService "github.com/vshevchenk0/bday-notifier/internal/service/user"
	"github.com/vshevchenk0/bday-notifier/pkg/jwt"
	"github.com/vshevchenk0/bday-notifier/pkg/logger"
)

type serviceProvider struct {
	config   *config.Config
	database *sqlx.DB

	tokenManager jwt.Manager
	logger       *slog.Logger

	userRepository         repository.UserRepository
	subscriptionRepository repository.SubscriptionRepository

	authService         service.AuthService
	userService         service.UserService
	subscriptionService service.SubscriptionService

	authMiddleware middleware.AuthMiddleware

	authHandler         *api.AuthHandler
	userHandler         *api.UserHandler
	subscriptionHandler *api.SubscriptionHandler
	router              http.Handler

	serverConfig *server.ServerConfig
}

func newServiceProvider(config *config.Config, db *sqlx.DB) *serviceProvider {
	s := &serviceProvider{}
	s.config = config
	s.database = db
	return s
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		panic("no config provided, service provider initialized wrong")
	}
	return s.config
}

func (s *serviceProvider) Database() *sqlx.DB {
	if s.database == nil {
		panic("no database provided, service provider initialized wrong")
	}
	return s.database
}

func (s *serviceProvider) TokenManager() jwt.Manager {
	if s.tokenManager == nil {
		managerConfig := &jwt.ManagerConfig{
			SigningKey: s.Config().JwtSigningKey,
			TokenTtl:   s.Config().JwtTokenTtl,
		}
		manager, err := jwt.NewManager(managerConfig)
		if err != nil {
			panic("failed to init token manager")
		}
		s.tokenManager = manager
	}
	return s.tokenManager
}

func (s *serviceProvider) Logger() *slog.Logger {
	if s.logger == nil {
		logger := logger.NewLogger(s.Config().Env)
		s.logger = logger
	}
	return s.logger
}

func (s *serviceProvider) UserRepository() repository.UserRepository {
	if s.userRepository == nil {
		s.userRepository = userRepository.NewRepository(s.Database())
	}
	return s.userRepository
}

func (s *serviceProvider) SubscriptionRepository() repository.SubscriptionRepository {
	if s.subscriptionRepository == nil {
		s.subscriptionRepository = subscriptionRepository.NewRepository(s.Database())
	}
	return s.subscriptionRepository
}

func (s *serviceProvider) AuthService() service.AuthService {
	if s.authService == nil {
		s.authService = authService.NewAuthService(
			s.UserRepository(),
			s.TokenManager(),
			s.Logger(),
		)
	}
	return s.authService
}

func (s *serviceProvider) UserService() service.UserService {
	if s.userService == nil {
		s.userService = userService.NewUserService(s.UserRepository(), s.Logger())
	}
	return s.userService
}

func (s *serviceProvider) SubscriptionService() service.SubscriptionService {
	if s.subscriptionService == nil {
		s.subscriptionService = subscriptionService.NewSubscriptionService(
			s.SubscriptionRepository(),
			s.Logger(),
		)
	}
	return s.subscriptionService
}

func (s *serviceProvider) AuthMiddleware() middleware.AuthMiddleware {
	if s.authMiddleware == nil {
		s.authMiddleware = middleware.NewAuthMiddleware(s.AuthService())
	}
	return s.authMiddleware
}

func (s *serviceProvider) AuthHandler() *api.AuthHandler {
	if s.authHandler == nil {
		s.authHandler = api.NewAuthHandler(s.AuthService())
	}
	return s.authHandler
}

func (s *serviceProvider) UserHandler() *api.UserHandler {
	if s.userHandler == nil {
		s.userHandler = api.NewUserHandler(s.UserService(), s.AuthMiddleware())
	}
	return s.userHandler
}

func (s *serviceProvider) SubscriptionHandler() *api.SubscriptionHandler {
	if s.subscriptionHandler == nil {
		s.subscriptionHandler = api.NewSubscriptionHandler(
			s.SubscriptionService(),
			s.AuthMiddleware(),
		)
	}
	return s.subscriptionHandler
}

func (s *serviceProvider) Router() http.Handler {
	if s.router == nil {
		s.router = api.NewRouter(s.AuthHandler(), s.SubscriptionHandler(), s.UserHandler())
	}
	return s.router
}

func (s *serviceProvider) ServerConfig() *server.ServerConfig {
	if s.serverConfig == nil {
		s.serverConfig = &server.ServerConfig{
			Host: s.Config().AppHost,
			Port: s.Config().AppPort,
		}
	}
	return s.serverConfig
}
