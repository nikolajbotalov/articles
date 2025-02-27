package handlers

import (
	"PersonalBlog/internal/handlers/article"
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRoutes(g *gin.Engine, uc usecases.ArticleUseCase, logger *zap.Logger) {
	articleRoutes := g.Group("/articles")
	{
		articleRoutes.GET("/all", handlers.GetArticles(uc, logger))
		articleRoutes.POST("/", handlers.CreateArticle(uc, logger))
		articleRoutes.GET("/:id", handlers.GetArticleByID(uc, logger))
		articleRoutes.PUT("/:id", handlers.UpdateArticleByID(uc, logger))
		articleRoutes.DELETE("/:id", handlers.DeleteArticleByID(uc, logger))
	}
}
