package usecases

import (
	"PersonalBlog/internal/domain"
	"context"
	"go.uber.org/zap"
)

func (uc *articleUseCase) GetAllArticles(ctx context.Context) ([]domain.Article, error) {
	uc.logger.Info("Fetching all articles")

	articles, err := uc.repo.GetAll(ctx)
	if err != nil {
		uc.logger.Error("Failed to fetch articles", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Successfully fetched articles", zap.Int("articles", len(articles)))
	return articles, nil
}
