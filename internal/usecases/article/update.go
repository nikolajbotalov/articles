package usecases

import (
	"PersonalBlog/internal/domain"
	"context"
	"go.uber.org/zap"
	"time"
)

func (uc *articleUseCase) UpdateArticle(ctx context.Context, id string, article *domain.Article) error {
	uc.logger.Info("Updating article", zap.String("id", id))

	article.UpdatedAt = time.Now()

	if err := uc.repo.Update(ctx, id, article); err != nil {
		uc.logger.Error("Failed to update article", zap.Error(err))
		return err
	}

	uc.logger.Info("Successfully updated article", zap.String("id", id))
	return nil
}
