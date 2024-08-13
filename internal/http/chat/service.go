package chat

import (
	"chat-application/internal/config"
	"chat-application/internal/storage/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"time"
)

type Service struct {
	log     *slog.Logger
	queries *postgres.Queries
}

func NewService(log *slog.Logger) *Service {
	const op = "chat.service.NewService"

	cfg := config.MustLoad()

	conn, err := pgx.Connect(context.Background(), fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.DB.DbUser, cfg.DB.DbPass, cfg.DB.DbHost, cfg.DB.DbPort, cfg.DB.DbName, cfg.DB.SslMode),
	)
	if err != nil {
		log.With("op", op).Warn("failed to connect to db")
	}

	queries := postgres.New(conn)

	return &Service{
		log:     log,
		queries: queries,
	}
}

func (s *Service) CreateUser(c context.Context, username, email, password string) error {
	ctx, cancel := context.WithTimeout(c, time.Duration(2)*time.Second)
	defer cancel()

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// TODO: logging

	err = s.queries.CreateUser(ctx, postgres.CreateUserParams{
		Username: username,
		Password: hashedPass,
		Email:    email,
	})
	if err != nil {
		return err
	}

	return nil
}
