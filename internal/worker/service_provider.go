package worker

import (
	"log/slog"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-notifier/internal/config"
	"github.com/vshevchenk0/bday-notifier/internal/repository"
	notificationRepository "github.com/vshevchenk0/bday-notifier/internal/repository/notification"
	"github.com/vshevchenk0/bday-notifier/internal/service"
	notificationService "github.com/vshevchenk0/bday-notifier/internal/service/notification"
	"github.com/vshevchenk0/bday-notifier/pkg/logger"
	"github.com/vshevchenk0/bday-notifier/pkg/mailer"
)

type serviceProvider struct {
	config   *config.Config
	database *sqlx.DB

	mailer mailer.Mailer
	logger *slog.Logger

	notificationRepository repository.NotificationRepository

	notificationService service.NotificationService
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

func (s *serviceProvider) Mailer() mailer.Mailer {
	if s.mailer == nil {
		mailerConfig := &mailer.MailerConfig{
			Email:           s.Config().MailerEmail,
			Password:        s.Config().MailerPassword,
			SmtpHost:        s.Config().MailerSmtpHost,
			SmtpPort:        s.Config().MailerSmtpPort,
			WaitBeforeRetry: s.Config().MailerWaitBeforeRetry,
			MaxRetriesCount: s.Config().MailerMaxRetriesCount,
			IncrementalWait: s.Config().MailerIncrementalWait,
		}
		s.mailer = mailer.NewMailer(mailerConfig, s.Logger())
	}
	return s.mailer
}

func (s *serviceProvider) Logger() *slog.Logger {
	if s.logger == nil {
		logger := logger.NewLogger(s.Config().Env)
		s.logger = logger
	}
	return s.logger
}

func (s *serviceProvider) NotificationRepository() repository.NotificationRepository {
	if s.notificationRepository == nil {
		s.notificationRepository = notificationRepository.NewRepository(s.Database())
	}
	return s.notificationRepository
}

func (s *serviceProvider) NotificationService() service.NotificationService {
	if s.notificationService == nil {
		s.notificationService = notificationService.NewNotificationService(
			s.NotificationRepository(),
			s.Mailer(),
			s.Logger(),
		)
	}
	return s.notificationService
}
