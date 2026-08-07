package app

import (
	"fmt"
	"lingva/services/analyzer/internal/adapters/out_docker"
	grpcapp "lingva/services/analyzer/internal/app/grpc"
	"lingva/services/analyzer/internal/core/service"
	"log/slog"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(log *slog.Logger, port int) *App {
	const op = "app.New"

	sandbox, err := out_docker.New(log)
	if err != nil {
		panic(fmt.Sprintf("%s: failed to initialize sandbox: %v", op, err))
	}
	analyzer := service.New(sandbox)

	grpcApp := grpcapp.New(log, port, analyzer)
	return &App{
		GRPCServer: grpcApp,
	}
}
