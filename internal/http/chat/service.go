package chat

import (
	"chat-application/internal/config"
	"chat-application/internal/storage/postgres"
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strconv"
	"time"
)

const (
	secretKey = "secret"
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

func (s *Service) Login(c context.Context, email, password string) (string, error) {
	log := s.log.With("op", "http.service.Login")

	ctx, cancel := context.WithTimeout(c, time.Duration(2)*time.Second)
	defer cancel()

	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		log.Error("Failed to get user", "error", err)
		return "", err
	}

	err = bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		log.Info("Invalid password")
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, struct {
		ID       string
		Username string
		jwt.RegisteredClaims
	}{
		ID:       strconv.Itoa(int(user.ID)),
		Username: user.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    strconv.Itoa(int(user.ID)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	})

	ss, err := token.SignedString([]byte(secretKey))
	if err != nil {
		log.Error("Failed to generate jwt token", "error", err)
		return "", err
	}

	return ss, nil
}
