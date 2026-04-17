package app

import (
	"assignments/simplebank/accounts"
	"assignments/simplebank/adapters/ginserver"
	"assignments/simplebank/adapters/monitoring"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"assignments/simplebank/adapters/postgres"

	"github.com/kelseyhightower/envconfig"
)

// App is the main application struct that holds all the dependencies and configurations.
var Version = "unknown"

type Config struct {
	HTTP     ginserver.Config
	Postgres postgres.Config

	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" default:"20s"`
}

func Start() error {
	var config Config
	envconfig.MustProcess("", &config)

	log := monitoring.Logger()
	log.Info("version " + Version)
	db, dbconn := postgres.Connect(config.Postgres)
	defer dbconn.Close()

	router, httpStart := ginserver.InitGin(config.HTTP)
	basePath := config.HTTP.Prefix
	httpGroup := router.Group(basePath)
	accounts.Initialize(db, httpGroup)

	_, stopHTTP := httpStart()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stopHTTP(ctx)
	return nil
}
