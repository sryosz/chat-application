package app

import (
	httpapp "chat-application/internal/app/http"
	"log/slog"
)

type App struct {
	HttpApp *httpapp.App
}

func New(log *slog.Logger) *App {
	httpApp := httpapp.NewApp(log)
	return &App{
		HttpApp: httpApp,
	}
}
