package httpapp

import (
	"context"
	"errors"
	"fmt"
	"lingva/api/gen"
	"lingva/services/gateway/internal/transport/rest"
	"lingva/services/gateway/internal/transport/websocket"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "lingva/docs"

	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	log        *slog.Logger
	httpServer *http.Server
	port       int
}

func New(log *slog.Logger, port int, runner gen.CodeRunServiceClient, analyzer gen.CodeAnalyzeServiceClient) *App {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	runHandler := websocket.NewCodeRunHandler(log, runner)
	analyzeHandler := rest.NewCodeAnalyzeHandler(log, analyzer)

	router.Get("/ws/run", runHandler)
	router.Post("/rest/analyze", analyzeHandler)
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("http://localhost:%d/swagger/doc.json", port))))
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	return &App{
		log:        log,
		httpServer: httpServer,
		port:       port,
	}
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	}
}

func (a *App) Run() error {
	const op = "httpapp.Run"

	log := a.log.With(slog.String("op", op))
	log.Info("gateway http server is running", slog.String("addr", a.httpServer.Addr))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (a *App) Stop() {
	const op = "httpapp.Stop"

	log := a.log.With(slog.String("op", op))
	log.Info("stopping gateway http server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.Shutdown(ctx); err != nil {
		log.Error(fmt.Sprintf("%s: failed to stop http server: %s", op, err.Error()))
	}
}
