package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"websocket/internal/controller" 
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port int
	clients *controller.ClientController 
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port: port,
		clients: controller.NewClientController(),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
