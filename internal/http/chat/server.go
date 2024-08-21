package chat

import (
	"chat-application/internal/models"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log/slog"
	"net/http"
)

type chatService interface {
	CreateUser(ctx context.Context, username, email, password string) error
	Login(c context.Context, email, password string) (string, error)
}

type Server struct {
	log *slog.Logger
	Hub *models.Hub
	chatService
	*gin.Engine
}

func NewServer(log *slog.Logger) *Server {
	server := gin.Default()
	service := NewService(log)

	return &Server{
		log:         log,
		chatService: service,
		Hub:         models.NewHub(),
		Engine:      server,
	}
}

func (s *Server) RegisterRoutes() {
	s.POST("/signup", s.createUser)
	s.POST("/login", s.login)
	s.POST("/logout", s.logout)

	s.POST("/ws/createRoom", s.createRoom)
	s.GET("/ws/joinRoom/:roomId", s.joinRoom)
	s.GET("/ws/getRooms", s.getRooms)
	s.GET("/ws/getClients/:roomId", s.getClients)
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

func (s *Server) createRoom(ctx *gin.Context) {
	reqBody := struct {
		ID   string
		Name string
	}{}

	if err := ctx.ShouldBindJSON(&reqBody); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	s.Hub.Rooms[reqBody.ID] = &models.Room{
		ID:      reqBody.ID,
		Name:    reqBody.Name,
		Clients: make(map[string]*models.Client),
	}

	ctx.JSON(http.StatusOK, reqBody)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) joinRoom(ctx *gin.Context) {
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	roomID := ctx.Param("roomId")
	clientID := ctx.Query("userId")
	username := ctx.Query("username")

	client := &models.Client{
		Conn:     conn,
		Message:  make(chan *models.Message, 10),
		ID:       clientID,
		RoomID:   roomID,
		Username: username,
	}

	msg := &models.Message{
		Content:  "A new user has joined the room",
		RoomID:   roomID,
		Username: username,
	}

	s.Hub.Register <- client
	s.Hub.Broadcast <- msg

	go client.WriteMessage()
	client.ReadMessage(s.Hub)
}

type RoomRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (s *Server) getRooms(ctx *gin.Context) {

	rooms := make([]RoomRes, 0)

	for _, r := range s.Hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:   r.ID,
			Name: r.Name,
		})
	}

	ctx.JSON(http.StatusOK, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (s *Server) getClients(ctx *gin.Context) {
	var clients []ClientRes

	roomId := ctx.Param("roomId")

	if _, ok := s.Hub.Rooms[roomId]; !ok {
		clients = make([]ClientRes, 0)
		ctx.JSON(http.StatusOK, clients)
	}

	for _, c := range s.Hub.Rooms[roomId].Clients {
		clients = append(clients, ClientRes{
			ID:       c.ID,
			Username: c.Username,
		})
	}

	ctx.JSON(http.StatusOK, clients)
}
