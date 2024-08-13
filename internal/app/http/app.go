package httpapp

import (
	"chat-application/internal/http/chat"
	"log/slog"
)

type App struct {
	httpServer *chat.Server
	log        *slog.Logger
}

func NewApp(log *slog.Logger) *App {
	server := chat.NewServer(log)
	server.RegisterRoutes()

	return &App{
		httpServer: server,
		log:        log,
	}
}

func (a *App) Run(addr string) {
	const op = "app.http.Run"

	log := a.log.With("op", op)

	log.Info("Starting HTTP server")

	go func() {
		err := a.httpServer.Run(addr)
		if err != nil {
			panic(err)
		}
	}()

	log.Info("HTTP server is running", slog.String("addr", addr))
}
