package worker

import (
	"context"
	"log/slog"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/vshevchenk0/bday-notifier/internal/config"
)

type Worker struct {
	serviceProdider *serviceProvider
}

func NewWorker(config *config.Config, db *sqlx.DB) (*Worker, error) {
	w := &Worker{}
	if err := w.initServiceProvider(config, db); err != nil {
		return nil, err
	}

	return w, nil
}

func (w *Worker) initServiceProvider(config *config.Config, db *sqlx.DB) error {
	w.serviceProdider = newServiceProvider(config, db)
	return nil
}

func (w *Worker) Run(ctx context.Context, doneChan chan struct{}) error {
	currentTime := time.Now()
	timeUntilNextDay := time.Until(currentTime.Add(time.Hour * time.Duration(24-currentTime.Hour())).Round(time.Hour))
	timeoutCtx, cancel := context.WithTimeout(ctx, timeUntilNextDay)
	defer cancel()
	if err := w.serviceProdider.NotificationService().NotifyUsers(timeoutCtx); err != nil {
		w.serviceProdider.Logger().Error("worker error", slog.String("error", err.Error()))
		return err
	}
	doneChan <- struct{}{}
	return nil
}
