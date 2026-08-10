package main

import (
	cfgLoader "lingva/pkg/config"
	"lingva/pkg/logger"
	"lingva/services/gateway/internal/app"
	"lingva/services/gateway/internal/config"
	"os"
	"os/signal"
)

// @title           Lingva Gateway API
// @host            localhost:8080
// @BasePath        /
func main() {
	var cfg config.Config
	cfgLoader.MustLoad("configs/gateway.yml", &cfg)

	log := logger.Setup()
	log.Info("starting gateway")

	application := app.New(log, cfg.Server.Port, cfg.GRPCClients.Runner, cfg.GRPCClients.Analyzer)
	go application.HTTPApp.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info("stopping gateway")
	application.HTTPApp.Stop()
	log.Info("gateway application stopped")
}
