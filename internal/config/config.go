package config

import (
	"time"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env string `env:"ENV"`

	AppHost string `env:"APP_HOST"`
	AppPort string `env:"APP_PORT"`

	JwtSigningKey string        `env:"JWT_SIGNING_KEY,unset"`
	JwtTokenTtl   time.Duration `env:"JWT_TOKEN_TTL" envDefault:"1h"`

	MailerEmail           string        `env:"MAILER_EMAIL"`
	MailerPassword        string        `env:"MAILER_PASSWORD"`
	MailerSmtpHost        string        `env:"MAILER_SMTP_HOST"`
	MailerSmtpPort        string        `env:"MAILER_SMTP_PORT"`
	MailerWaitBeforeRetry time.Duration `env:"MAILER_WAIT_BEFORE_RETRY"`
	MailerMaxRetriesCount int           `env:"MAILER_MAX_RETRIES_COUNT"`
	MailerIncrementalWait bool          `env:"MAILER_INCREMENTAL_WAIT"`

	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     int    `env:"DB_PORT"`
	DatabaseName     string `env:"DB_NAME"`
}

func MustLoad() *Config {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		panic("failed to load config")
	}
	return cfg
}
