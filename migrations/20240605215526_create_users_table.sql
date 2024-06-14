-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
	id uuid primary key default gen_random_uuid(),
	email varchar(255) not null unique,
	password_hash varchar(255) not null,
	name varchar(255) not null,
	surname varchar(255) not null,
	birthday_date date not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
-- +goose StatementEnd
