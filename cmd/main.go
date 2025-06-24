package main

import (
	"fmt"

	"github.com/DblMOKRQ/workmate_test_task/internal/config"
	"github.com/DblMOKRQ/workmate_test_task/internal/router"
	"github.com/DblMOKRQ/workmate_test_task/internal/router/handlers"
	"github.com/DblMOKRQ/workmate_test_task/internal/storage"
	logger "github.com/DblMOKRQ/workmate_test_task/pkg"
	"go.uber.org/zap"
)

func main() {

	cfg := config.MustLoad()
	log, err := logger.NewLogger()

	if err != nil {
		panic(err)
	}

	store := storage.NewTaskStorage()
	handler := handlers.NewHandler(store, log)
	rt := router.NewRouter(handler, log)

	addr := fmt.Sprintf("%s:%d", cfg.HTTP.Host, cfg.HTTP.Port)
	if err := rt.Run(addr); err != nil {
		log.Fatal("Failed to start server", zap.Error(err))
	}

	// Ожидание сигналов и graceful shutdown
	if err := rt.GracefulShutdown(); err != nil {
		log.Error("Shutdown error", zap.Error(err))
	}
}
