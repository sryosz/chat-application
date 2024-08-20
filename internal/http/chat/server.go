package chat

import (
	"context"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type chatService interface {
	CreateUser(ctx context.Context, username, email, password string) error
	Login(c context.Context, email, password string) (string, error)
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
	s.POST("/login", s.login)
	s.POST("/logout", s.logout)
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
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}

	ctx.JSON(http.StatusCreated, true)
}

func (s *Server) login(ctx *gin.Context) {
	reqBody := struct {
		Email    string
		Password string
	}{}

	err := ctx.BindJSON(&reqBody)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	token, err := s.chatService.Login(ctx, reqBody.Email, reqBody.Password)
	if err != nil {
		ctx.AbortWithStatus(http.StatusConflict)
		return
	}

	ctx.SetCookie("jwt", token, 3600, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"message": "authorized"})
}

func (s *Server) logout(ctx *gin.Context) {
	ctx.SetCookie("jwt", "", -1, "", "", false, true)
	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
