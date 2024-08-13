package chat

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type chatService interface {
	CreateUser(ctx context.Context, username, email, password string) error
}

type Server struct {
	log *slog.Logger
	chatService
	*gin.Engine
}

func NewServer(log *slog.Logger) *Server {
	server := gin.Default()
	service := NewService(log)

	return &Server{
		log:         log,
		chatService: service,
		Engine:      server,
	}
}

func (s *Server) RegisterRoutes() {
	s.POST("/signup", s.createUser)
}

func (s *Server) createUser(ctx *gin.Context) {
	reqBody := struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	err = s.chatService.CreateUser(ctx, reqBody.Username, reqBody.Email, reqBody.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusNotFound)
		return
	}

	ctx.JSON(http.StatusCreated, true)
}
