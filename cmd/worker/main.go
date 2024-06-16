package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"
	"github.com/vshevchenk0/bday-greeter/internal/config"
	"github.com/vshevchenk0/bday-greeter/internal/worker"
	"github.com/vshevchenk0/bday-greeter/pkg/postgresql"
)

func main() {
	config := config.MustLoad()

	dbConfig := &postgresql.PostgresqlConfig{
		User:         config.DatabaseUser,
		Password:     config.DatabasePassword,
		Host:         config.DatabaseHost,
		Port:         config.DatabasePort,
		DatabaseName: config.DatabaseName,
	}
	db, err := postgresql.NewPostgresqlDB(dbConfig)
	if err != nil {
		panic("failed to connect to database")
	}
	defer db.Close()

	w, err := worker.NewWorker(config, db)
	if err != nil {
		panic(fmt.Errorf("failed to initialize worker: %v", err))
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)

	doneChan := make(chan struct{}, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		if err := w.Run(ctx, doneChan); err != nil {
			return
		}
	}()

	select {
	case <-shutdownChan:
		cancel()
		<-doneChan
		return
	case <-doneChan:
		return
	}
}
