package app

import (
	"fmt"
	"lingva/services/runner/internal/adapters/out_docker"
	grpcapp "lingva/services/runner/internal/app/grpc"
	"lingva/services/runner/internal/core/service"
	"log/slog"
)

type App struct {
	GRPCApp *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	const op = "App.New"

	sandbox, err := out_docker.New(log)
	if err != nil {
		panic(fmt.Sprintf("%s: cannot init docker client: %s", op, err))
	}
	runner := service.New(sandbox)

	grpcApp := grpcapp.New(log, port, runner)
	return &App{
		GRPCApp: grpcApp,
	}
}
