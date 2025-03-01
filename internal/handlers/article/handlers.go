package handlers

import (
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupArticleRoutes(g *gin.Engine, uc usecases.ArticleUseCase, logger *zap.Logger) {
	articleRoutes := g.Group("/articles")
	{
		articleRoutes.GET("/all", GetArticles(uc, logger))
		articleRoutes.POST("/", CreateArticle(uc, logger))
		articleRoutes.GET("/:id", GetArticleByID(uc, logger))
		articleRoutes.PUT("/:id", UpdateArticleByID(uc, logger))
		articleRoutes.DELETE("/:id", DeleteArticleByID(uc, logger))
	}
}
