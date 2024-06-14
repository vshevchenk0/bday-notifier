package app

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-greeter/internal/config"
	"github.com/vshevchenk0/bday-greeter/internal/server"
)

type App struct {
	serviceProdider *serviceProvider
	server          *server.Server
}

func NewApp(config *config.Config, db *sqlx.DB) (*App, error) {
	a := &App{}
	if err := a.initServiceProvider(config, db); err != nil {
		return nil, err
	}

	if err := a.initServer(); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initServiceProvider(config *config.Config, db *sqlx.DB) error {
	a.serviceProdider = newServiceProvider(config, db)
	return nil
}

func (a *App) initServer() error {
	a.server = server.NewServer(
		a.serviceProdider.ServerConfig(),
		a.serviceProdider.Router(),
		a.serviceProdider.Logger(),
	)
	return nil
}

func (a *App) Run() {
	a.server.Start()
}

func (a *App) Stop(ctx context.Context) error {
	return a.server.Stop(ctx)
}
