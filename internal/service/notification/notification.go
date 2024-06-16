package notification

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/vshevchenk0/bday-greeter/internal/repository"
	"github.com/vshevchenk0/bday-greeter/pkg/mailer"
)

type notificationService struct {
	notificationRepository repository.NotificationRepository
	mailer                 mailer.Mailer
	logger                 *slog.Logger
}

func NewNotificationService(
	notificationRepository repository.NotificationRepository,
	mailer mailer.Mailer,
	logger *slog.Logger,
) *notificationService {
	return &notificationService{
		notificationRepository: notificationRepository,
		mailer:                 mailer,
		logger:                 logger,
	}
}

func (s *notificationService) NotifyUsers(ctx context.Context) error {
	type UserId string
	type UserInfo struct {
		name              string
		surname           string
		birthdayDate      time.Time
		daysUntilBirthday int
		subscribersEmails []string
	}

	tx, err := s.notificationRepository.GetLock(ctx)
	defer func() {
		err := tx.Commit()
		if err != nil {
			s.logger.Error("error committing transaction", slog.String("error", err.Error()))
		}
	}()

	if errors.Is(err, repository.ErrLockTaken) {
		s.logger.Info("job is already done by other worker")
		return nil
	}
	if err != nil {
		s.logger.Error("failed to take lock", slog.String("error", err.Error()))
		return err
	}

	notificationRecords, err := s.notificationRepository.FindUsersToNotify(ctx, tx)
	if err != nil {
		s.logger.Error("failed to retrieve notification records", slog.String("error", err.Error()))
		_ = tx.Rollback()
		return err
	}

	notificationsMap := make(map[UserId]UserInfo)
	for _, v := range notificationRecords {
		userInfo, ok := notificationsMap[UserId(v.BirthdayUserId)]
		if !ok {
			notificationsMap[UserId(v.BirthdayUserId)] = UserInfo{
				name:              v.BirthdayUserName,
				surname:           v.BirthdayUserSurname,
				birthdayDate:      v.BirthdayDate,
				daysUntilBirthday: v.DaysUntilBirthday,
				subscribersEmails: []string{v.SubscriberEmail},
			}
		} else {
			userInfo.subscribersEmails = append(userInfo.subscribersEmails, v.SubscriberEmail)
			notificationsMap[UserId(v.BirthdayUserId)] = userInfo
		}
	}

	wg := &sync.WaitGroup{}
	for _, userInfo := range notificationsMap {
		wg.Add(1)
		subject := fmt.Sprintf("Birthday of %s %s", userInfo.name, userInfo.surname)
		body := fmt.Sprintf(
			"%s %s will celebrate birthday in %d days, on %d of %s!",
			userInfo.name,
			userInfo.surname,
			userInfo.daysUntilBirthday,
			userInfo.birthdayDate.Day(),
			time.Month(userInfo.birthdayDate.Month()).String(),
		)
		go s.mailer.Send(ctx, wg, userInfo.subscribersEmails, subject, body)
	}
	wg.Wait()
	return nil
}
