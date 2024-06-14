package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/vshevchenk0/bday-greeter/internal/config"
	"github.com/vshevchenk0/bday-greeter/pkg/mailer"
)

func main() {
	config := config.MustLoad()
	fmt.Println(config.MailerEmail, config.MailerPassword)
	mailer := mailer.NewMailer(
		&mailer.MailerConfig{
			Email:    config.MailerEmail,
			Password: config.MailerPassword,
			SmtpHost: config.MailerSmtpHost,
			SmtpPort: config.MailerSmtpPort,
		},
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo})),
	)
	err := mailer.SendEmail([]string{"sepabe3943@eqvox.com"}, "test", "test")
	if err != nil {
		panic(err)
	}
}
