package handlers

import (
	"PersonalBlog/internal/handlers/article"
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(g *gin.Engine, uc usecases.ArticleUseCase) {
	articleRoutes := g.Group("/articles")
	{
		articleRoutes.GET("/all", handlers.GetArticles(uc))
		articleRoutes.POST("/", handlers.CreateArticle(uc))
		articleRoutes.GET("/:id", handlers.GetArticleByID(uc))
		articleRoutes.PUT("/:id", handlers.UpdateArticleByID(uc))
		articleRoutes.DELETE("/:id", handlers.DeleteArticleByID(uc))
	}
}
