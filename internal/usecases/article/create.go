package usecases

import (
	"PersonalBlog/internal/domain"
	"context"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

func (uc *articleUseCase) CreateArticle(ctx context.Context, article *domain.Article) error {
	uc.logger.Info("Creating article", zap.String("title", article.Title))

	article.ID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	if err := uc.repo.Create(ctx, article); err != nil {
		uc.logger.Error("Failed to create article", zap.Error(err))
		return err
	}

	uc.logger.Info("Successfully created article", zap.String("title", article.ID))
	return nil
}
