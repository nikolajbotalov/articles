package handlers

import (
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// GetArticles godoc
// @Summary Получить все статьи
// @Description Возвращает список всех статей
// @Tags articles
// @Accept  json
// @Produce  json
// @Success 200 {array} domain.Article
// @Failure 500 {object} map[string]string
// @Router /articles/all [get]
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
