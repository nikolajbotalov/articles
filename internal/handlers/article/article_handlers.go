package handlers

import (
	"PersonalBlog/internal/domain"
	usecases "PersonalBlog/internal/usecases/article"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func GetArticles(uc usecases.ArticleUseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		articles, err := uc.GetAllArticles(ctx)
		if err != nil {
			logger.Error("Failed to get articles", zap.Error(err))
			handleError(c, http.StatusInternalServerError, "Failed to get articles", err, logger)
		}

		c.JSON(http.StatusOK, articles)
	}
}

func CreateArticle(uc usecases.ArticleUseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var newArticle domain.Article

		// Парсим JSON из тела запроса
		if err := parseJSON(c, &newArticle, logger); err != nil {
			return
		}

		// Валидация данных
		if err := validateArticle(c, newArticle, logger); err != nil {
			return
		}

		if err := uc.CreateArticle(ctx, &newArticle); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to create article", err, logger)
			return
		}

		// Возвращаем созданную статью
		c.JSON(http.StatusCreated, newArticle)
	}
}

func GetArticleByID(uc usecases.ArticleUseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		article, err := uc.GetArticleByID(ctx, id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to get article", err, logger)
		}

		// Возвращаем статью в формате JSON
		c.JSON(http.StatusOK, article)
	}
}

func UpdateArticleByID(uc usecases.ArticleUseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		// Парсим JSON из тела запроса
		var updateArticle domain.Article
		if err := parseJSON(c, &updateArticle, logger); err != nil {
			return
		}

		// Валидация данных
		if err := validateArticle(c, updateArticle, logger); err != nil {
			return
		}

		err := uc.UpdateArticle(ctx, id, &updateArticle)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update article", err, logger)
		}

		// Возвращаем обновленную статью
		updateArticle.ID = id
		c.JSON(http.StatusOK, updateArticle)
	}
}

func DeleteArticleByID(uc usecases.ArticleUseCase, logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		err := uc.DeleteArticle(ctx, id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete article", err, logger)
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
	}
}

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
