-- +goose Up
-- +goose StatementBegin
CREATE TABLE subscriptions (
	user_id uuid references users (id) on delete cascade,
	subscriber_id uuid references users (id) on delete cascade,
	notify_before_days integer not null default 0,
	unique (user_id, subscriber_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE subscriptions;
-- +goose StatementEnd
