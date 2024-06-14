package user

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/vshevchenk0/bday-greeter/internal/model"
	"github.com/vshevchenk0/bday-greeter/internal/repository"
)

type userRepository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *userRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) CreateUser(
	ctx context.Context, email, name, surname, passwordHash string, birthdayDate time.Time,
) (string, error) {
	var id string
	query := `
		INSERT INTO users (email, password_hash, name, surname, birthday_date)
		VALUES ($1, $2, $3, $4, $5) RETURNING id;
	`
	row := r.db.QueryRowxContext(ctx, query, email, passwordHash, name, surname, birthdayDate)
	err := row.Scan(&id)
	if err, ok := err.(*pq.Error); ok {
		// check unique constraint violation
		if err.Code == "23505" {
			return "", repository.ErrEmailIsNotUnique
		}
	}
	return id, err
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (model.User, error) {
	var user model.User
	query := "SELECT * FROM users WHERE email=$1;"
	err := r.db.GetContext(ctx, &user, query, email)
	return user, err
}

func (r *userRepository) FindAllUsers(ctx context.Context, userId string) ([]model.User, error) {
	var users []model.User
	// ::text cast if for proper date display
	query := "SELECT id, name, surname, birthday_date::text FROM users WHERE id != $1;"
	err := r.db.SelectContext(ctx, &users, query, userId)
	return users, err
}

func (r *userRepository) FindUsersSubscribedTo(ctx context.Context, userId string) ([]model.User, error) {
	var users []model.User
	query := `
		SELECT id, name, surname, birthday_date FROM users
		WHERE id in (
			SELECT user_id FROM subscriptions
			WHERE subscriber_id = $1
		);
	`
	err := r.db.SelectContext(ctx, &users, query, userId)
	return users, err
}
