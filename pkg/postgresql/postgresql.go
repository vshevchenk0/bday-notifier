package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PostgresqlConfig struct {
	User         string
	Password     string
	Host         string
	Port         int
	DatabaseName string
}

func NewPostgresqlDB(config *PostgresqlConfig) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User, config.Password, config.Host, config.Port, config.DatabaseName,
	)
	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return db, nil
}
