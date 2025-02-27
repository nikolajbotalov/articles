package handlers

import (
	"PersonalBlog/internal/domain"
	usecases "PersonalBlog/internal/usecases/article"
	"fmt"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

func GetArticles(uc usecases.ArticleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		articles, err := uc.GetAllArticles(ctx)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to get articles", err)
		}

		c.JSON(http.StatusOK, articles)
	}
}

func CreateArticle(uc usecases.ArticleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var newArticle domain.Article

		// Парсим JSON из тела запроса
		if err := parseJSON(c, &newArticle); err != nil {
			return
		}

		// Валидация данных
		if err := validateArticle(c, newArticle); err != nil {
			return
		}

		if err := uc.CreateArticle(ctx, &newArticle); err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to create article", err)
			return
		}

		// Возвращаем созданную статью
		c.JSON(http.StatusCreated, newArticle)
	}
}

func GetArticleByID(uc usecases.ArticleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		article, err := uc.GetArticleByID(ctx, id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to get article", err)
		}

		// Возвращаем статью в формате JSON
		c.JSON(http.StatusOK, article)
	}
}

func UpdateArticleByID(uc usecases.ArticleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		// Парсим JSON из тела запроса
		var updateArticle domain.Article
		if err := parseJSON(c, &updateArticle); err != nil {
			return
		}

		// Валидация данных
		if err := validateArticle(c, updateArticle); err != nil {
			return
		}

		err := uc.UpdateArticle(ctx, id, &updateArticle)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to update article", err)
		}

		// Возвращаем обновленную статью
		updateArticle.ID = id
		c.JSON(http.StatusOK, updateArticle)
	}
}

func DeleteArticleByID(uc usecases.ArticleUseCase) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		// Получаем ID статьи из параметров запроса
		id := c.Param("id")

		err := uc.DeleteArticle(ctx, id)
		if err != nil {
			handleError(c, http.StatusInternalServerError, "Failed to delete article", err)
		}

		// Возвращаем успешный ответ
		c.JSON(http.StatusOK, gin.H{"message": "Article deleted successfully"})
	}
}

// парсит JSON из тела запроса и возвращает ошибку, если что-то пошло не так
func parseJSON(c *gin.Context, article *domain.Article) error {
	if err := c.BindJSON(article); err != nil {
		slog.Error("Failed to bind JSON", "error", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return err
	}

	return nil
}

// проверка на пустые данные в полях статьи
func validateArticle(c *gin.Context, article domain.Article) error {
	if article.Author == "" || article.Title == "" || article.Content == "" {
		slog.Error("Validation failed", "article", article)
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Author, title, and content are required"})
		return fmt.Errorf("validation failed")
	}

	return nil
}

func handleError(c *gin.Context, statusCode int, message string, err error) {
	if err != nil {
		slog.Error(message, "error", err)
	} else {
		slog.Error(message)
	}
	c.AbortWithStatusJSON(statusCode, gin.H{"error": message})
}
