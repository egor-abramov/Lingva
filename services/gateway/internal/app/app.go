package app

import (
	"fmt"
	httpapp "lingva/services/gateway/internal/app/http"
	"lingva/services/gateway/internal/clients"
	"log/slog"
)

type App struct {
	HTTPApp *httpapp.App
}

func New(log *slog.Logger, port int, runnerTarget, analyzerTarget string) *App {
	const op = "App.New"

	runnerClient, err := clients.NewRunner(runnerTarget)
	if err != nil {
		panic(fmt.Sprintf("%s: cannot connect to runner service: %s", op, err))
	}

	analyzerClient, err := clients.NewAnalyzer(analyzerTarget)
	if err != nil {
		log.Error(fmt.Sprintf("%s: cannot connect to analyzer service: %s", op, err))
	}

	httpApp := httpapp.New(log, port, runnerClient, analyzerClient)
	return &App{
		HTTPApp: httpApp,
	}
}
