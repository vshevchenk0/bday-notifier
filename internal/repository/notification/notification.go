package notification

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-greeter/internal/model"
	"github.com/vshevchenk0/bday-greeter/internal/repository"
)

type notificationRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *notificationRepository {
	return &notificationRepository{
		db: db,
	}
}

func (r *notificationRepository) GetLock(ctx context.Context) (*sqlx.Tx, error) {
	opts := &sql.TxOptions{
		Isolation: 0,
		ReadOnly:  true,
	}
	tx, err := r.db.BeginTxx(ctx, opts)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, "LOCK TABLE ONLY subscriptions IN SHARE UPDATE EXCLUSIVE MODE NOWAIT;")
	if err != nil {
		return nil, repository.ErrLockTaken
	}
	return tx, nil
}

func (r *notificationRepository) FindUsersToNotify(ctx context.Context, tx *sqlx.Tx) ([]model.Notification, error) {
	var notifications []model.Notification
	query := `
		SELECT u2.email subscriber_email, u1.id birthday_user_id, u1.name birthday_user_name,
		u1.surname birthday_user_surname, u1.birthday_date birthday_date, s.notify_before_days days_until_birthday
		FROM subscriptions s
		LEFT JOIN users u1 on u1.id = s.user_id
		LEFT JOIN users u2 on u2.id = s.subscriber_id
		WHERE DATE_PART('day', u1.birthday_date) >= DATE_PART('day', CURRENT_DATE)
		AND DATE_PART('day', u1.birthday_date) <= DATE_PART('day', CURRENT_DATE) + 7
		AND DATE_PART('day', u1.birthday_date) = DATE_PART('day', CURRENT_DATE) + s.notify_before_days
		AND DATE_PART('month', u1.birthday_date) = DATE_PART('month', CURRENT_DATE);
	`
	err := tx.SelectContext(ctx, &notifications, query)
	return notifications, err
}
