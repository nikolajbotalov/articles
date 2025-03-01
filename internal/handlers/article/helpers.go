package handlers

import (
	"PersonalBlog/internal/domain"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// парсит JSON из тела запроса и возвращает ошибку, если что-то пошло не так
func parseJSON(c *gin.Context, article *domain.Article, logger *zap.Logger) error {
	if err := c.BindJSON(article); err != nil {
		logger.Error("Failed to bind JSON", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return err
	}

	return nil
}

// проверка на пустые данные в полях статьи
func validateArticle(c *gin.Context, article domain.Article, logger *zap.Logger) error {
	if article.Author == "" || article.Title == "" || article.Content == "" {
		logger.Error("Validation failed", zap.Any("article", article))
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Author, title, and content are required"})
		return fmt.Errorf("validation failed")
	}

	return nil
}

func handleError(c *gin.Context, statusCode int, message string, err error, logger *zap.Logger) {
	if err != nil {
		logger.Error(message, zap.Error(err))
	} else {
		logger.Error(message)
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
