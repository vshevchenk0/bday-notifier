package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/vshevchenk0/bday-greeter/internal/app"
	"github.com/vshevchenk0/bday-greeter/internal/config"
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

	a, err := app.NewApp(config, db)
	if err != nil {
		panic(fmt.Errorf("failed to initialize app: %v", err))
	}

	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, syscall.SIGINT, syscall.SIGTERM)
	a.Run()

	<-shutdownChan
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := a.Stop(ctx); err != nil {
		fmt.Printf("failed to gracefully shutdown")
	}
}
