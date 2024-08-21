package main

import (
	"chat-application/internal/app"
	"chat-application/internal/config"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log := setupLogger()

	log.Info("Starting application")

	application := app.New(log)

	cfg := config.MustLoad()

	application.HttpApp.Run(fmt.Sprintf(":%d", cfg.Http.Port))

	log.Info("Application started")

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	sign := <-stop

	log.Info("stopping:", slog.String("signal", sign.String()))

	log.Info("Stopping service")
}

func setupLogger() *slog.Logger {
	return slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
}
