package router

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DblMOKRQ/workmate_test_task/internal/router/handlers"
	"github.com/DblMOKRQ/workmate_test_task/internal/router/middleware"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Router struct {
	router  *mux.Router
	handler *handlers.Handler
	log     *zap.Logger
	server  *http.Server
}

func NewRouter(handler *handlers.Handler, log *zap.Logger) *Router {
	r := mux.NewRouter()
	r.HandleFunc("/create_task", handler.CreateTask).Methods(http.MethodPost)
	r.HandleFunc("/tasks/{id}", handler.GetTask).Methods(http.MethodGet)
	r.HandleFunc("/tasks", handler.DelteTask).Methods(http.MethodDelete)

	return &Router{
		router:  r,
		handler: handler,
		log:     log,
	}
}

func (r *Router) Run(addr string) error {
	middleware := middleware.LoggingMiddleware(r.router, r.log)

	r.server = &http.Server{
		Addr:    addr,
		Handler: middleware,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		r.log.Info("Starting server", zap.String("address", addr))
		if err := r.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			r.log.Fatal("Server failed to start", zap.Error(err))
		}
	}()

	return nil
}

func (r *Router) GracefulShutdown() error {
	r.log.Info("Starting graceful shutdown...")

	// Создаем контекст с таймаутом для shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Ожидаем сигналы завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Пытаемся корректно завершить работу сервера
	if err := r.server.Shutdown(ctx); err != nil {
		r.log.Error("Server shutdown failed", zap.Error(err))
		return err
	}

	r.log.Info("Server gracefully stopped")
	return nil
}
