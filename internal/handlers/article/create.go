package handlers

import (
	"PersonalBlog/internal/domain"
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// CreateArticle godoc
// @Summary Создать статью
// @Description Создает новую статью
// @Tags articles
// @Accept  json
// @Produce  json
// @Param article body domain.Article true "Данные статьи"
// @Success 201 {object} domain.Article
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles [post]
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
