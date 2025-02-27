package main

import (
	"PersonalBlog/internal/adapters/db"
	"PersonalBlog/internal/config"
	"PersonalBlog/internal/handlers"
	repositories "PersonalBlog/internal/repositories/article"
	usecases "PersonalBlog/internal/usecases/article"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
)

func main() {
	cfg := config.GetConfig()

	dbPool, err := db.NewPostgreSQLDB(cfg.PostgreSQL)
	if err != nil {
		slog.Error("Failed to initialize database", err)
		return
	}
	defer dbPool.Close()

	articleRepo := repositories.NewArticleRepository(dbPool)
	articleUseCase := usecases.NewArticleUseCase(articleRepo)

	router := gin.Default()
	handlers.SetupRoutes(router, articleUseCase)

	address := fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	serverStartErr := router.Run(address)
	if serverStartErr != nil {
		slog.Info("Unable to start server", serverStartErr)
		return
	}
}
