package model

type Subscription struct {
	UserId           string `json:"user_id,omitempty" db:"user_id"`
	SubscriberId     string `json:"subscriber_id,omitempty" db:"subscriber_id"`
	NotifyBeforeDays int    `json:"notify_before_days,omitempty" db:"notify_before_days"`
}
