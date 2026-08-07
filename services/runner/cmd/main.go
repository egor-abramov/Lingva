package main

import (
	cfgLoader "lingva/pkg/config"
	"lingva/pkg/logger"
	"lingva/services/runner/internal/app"
	"lingva/services/runner/internal/config"
	"os"
	"os/signal"
)

func main() {
	var cfg config.Config
	cfgLoader.MustLoad("configs/runner.yml", &cfg)

	log := logger.Setup()
	log.Info("starting runner")

	application := app.New(log, cfg.GRPC.Port)

	go application.GRPCApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop
	log.Info("stopping runner application")
	application.GRPCApp.Stop()
	log.Info("runner application stopped")
}
