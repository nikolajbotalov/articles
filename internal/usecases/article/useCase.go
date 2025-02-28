package usecases

import (
	"PersonalBlog/internal/domain"
	articleRepo "PersonalBlog/internal/repositories/article"
	"context"
	"go.uber.org/zap"
)

type ArticleUseCase interface {
	GetAllArticles(ctx context.Context) ([]domain.Article, error)
	GetArticleByID(ctx context.Context, id string) (*domain.Article, error)
	CreateArticle(ctx context.Context, article *domain.Article) error
	UpdateArticle(ctx context.Context, id string, article *domain.Article) error
	DeleteArticle(ctx context.Context, id string) error
}

type articleUseCase struct {
	repo   articleRepo.Repository
	logger *zap.Logger
}

func NewArticleUseCase(repo articleRepo.Repository, logger *zap.Logger) ArticleUseCase {
	return &articleUseCase{
		repo:   repo,
		logger: logger,
	}
}
