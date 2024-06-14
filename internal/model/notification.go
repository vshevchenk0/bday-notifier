package model

import "time"

type Notification struct {
	BirthdayUserId      string    `db:"birthday_user_id"`
	BirthdayUserName    string    `db:"birthday_user_name"`
	BirthdayUserSurname string    `db:"birthday_user_surname"`
	BirthdayDate        time.Time `db:"birthday_date"`
	SubscriberEmail     string    `db:"subscriber_email"`
	DaysUntilBirthday   int       `db:"days_until_birthday"`
}
