package usecases

import (
	"PersonalBlog/internal/domain"
	"context"
	"go.uber.org/zap"
)

func (uc *articleUseCase) GetArticleByID(ctx context.Context, id string) (*domain.Article, error) {
	uc.logger.Info("Fetching article by ID", zap.String("id", id))

	article, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		uc.logger.Error("Failed to fetch article by ID", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Successfully fetched article by ID", zap.String("id", id))
	return article, nil
}
