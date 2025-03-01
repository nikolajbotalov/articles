package handlers

import (
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetArticleByID godoc
// @Summary Получить статью по ID
// @Description Возвращает статью по её ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param id path string true "ID статьи"
// @Success 200 {object} domain.Article
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [get]
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
