package handlers

import (
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

// DeleteArticleByID godoc
// @Summary Удалить статью по ID
// @Description Удаляет статью по её ID
// @Tags articles
// @Accept  json
// @Produce  json
// @Param id path string true "ID статьи"
// @Success 200 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /articles/{id} [delete]
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
