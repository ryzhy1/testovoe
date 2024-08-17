package app

import (
	server "awesomeProject/internal/app/app"
	"awesomeProject/internal/routes"
	"awesomeProject/internal/services"
	"awesomeProject/internal/storage/postgres"
	"github.com/gin-gonic/gin"
	"log/slog"
	"time"
)

type App struct {
	HTTPServer *server.Server
}

func New(log *slog.Logger, serverPort, storagePath string, tokenTTL time.Duration) *App {
	storage, err := postgres.NewPostgresDB(storagePath)
	if err != nil {
		panic(err)
	}

	tokenService := services.NewTokenService(log, storage)

	r := gin.Default()
	routes.InitRoutes(r, tokenService)

	server := server.NewServer(log, serverPort, r)
	if err != nil {
		panic(err)
	}

	return &App{
		HTTPServer: server,
	}
}
