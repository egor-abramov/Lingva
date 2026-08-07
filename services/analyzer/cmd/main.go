package main

import (
	cfgLoader "lingva/pkg/config"
	"lingva/pkg/logger"
	"lingva/services/analyzer/internal/app"
	"lingva/services/analyzer/internal/config"
	"os"
	"os/signal"
)

func main() {
	var cfg config.Config
	cfgLoader.MustLoad("configs/analyzer.yml", &cfg)

	log := logger.Setup()
	log.Info("starting analyzer")

	application := app.New(log, cfg.GRPC.Port)
	go application.GRPCServer.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info("stopping analyzer application")
	application.GRPCServer.Stop()
	log.Info("analyzer application stopped")
}
