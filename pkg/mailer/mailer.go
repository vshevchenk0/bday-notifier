package mailer

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/smtp"
	"sync"
	"time"
)

type MailerConfig struct {
	Email           string
	Password        string
	SmtpHost        string
	SmtpPort        string
	WaitBeforeRetry time.Duration
	MaxRetriesCount int
	IncrementalWait bool
}

type Mailer interface {
	Send(ctx context.Context, wg *sync.WaitGroup, addresses []string, subject, body string)
}

type mailer struct {
	auth            smtp.Auth
	email           string
	smtpAddr        string
	waitBeforeRetry time.Duration
	maxRetriesCount int
	incrementalWait bool
	logger          *slog.Logger
}

func NewMailer(config *MailerConfig, logger *slog.Logger) *mailer {
	return &mailer{
		auth:            smtp.PlainAuth("", config.Email, config.Password, config.SmtpHost),
		email:           config.Email,
		smtpAddr:        net.JoinHostPort(config.SmtpHost, config.SmtpPort),
		waitBeforeRetry: config.WaitBeforeRetry,
		maxRetriesCount: config.MaxRetriesCount,
		incrementalWait: config.IncrementalWait,
		logger:          logger,
	}
}

func (m *mailer) sendEmail(addresses []string, subject string, body string) error {
	formattedSubject := fmt.Sprintf("Subject: %s\n", subject)
	message := []byte(formattedSubject + body)
	err := smtp.SendMail(m.smtpAddr, m.auth, m.email, addresses, message)
	return err
}

func (m *mailer) Send(ctx context.Context, wg *sync.WaitGroup, addresses []string, subject, body string) {
	defer wg.Done()
	queue := make(chan struct{}, 1)
	defer close(queue)
	queue <- struct{}{}
	retriesCount := 0

	for range queue {
		err := m.sendEmail(addresses, subject, body)
		if err == nil {
			m.logger.Info("successfully sent emails", slog.String("subject", subject))
			return
		}

		// requeue send message job when error occurs
		m.logger.Error("error sending email", slog.String("error", err.Error()))
		if retriesCount >= m.maxRetriesCount {
			m.logger.Warn("max retries reached, emails were not sent")
			return
		}
		retriesCount++

		var waitBeforeRetry time.Duration
		if m.incrementalWait {
			waitBeforeRetry = m.waitBeforeRetry * time.Duration(retriesCount)
		} else {
			waitBeforeRetry = m.waitBeforeRetry
		}

		m.logger.Info(fmt.Sprintf("retrying in %s", waitBeforeRetry))
		select {
		case <-time.After(waitBeforeRetry):
			queue <- struct{}{}
		case <-ctx.Done():
			m.logger.Warn("execution context was closed, exiting")
			return
		}
	}
}
