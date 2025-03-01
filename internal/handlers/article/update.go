package handlers

import (
	"PersonalBlog/internal/domain"
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// UpdateArticleByID godoc
// @Summary Обновить статью по ID
// @Description Обновляет статью по её ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param id path string true "ID статьи"
// @Param article body domain.Article true "Обновленные данные статьи"
// @Success 200 {object} domain.Article
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [put]
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
