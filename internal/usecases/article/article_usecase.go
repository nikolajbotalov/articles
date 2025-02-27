package usecases

import (
	"PersonalBlog/internal/domain"
	articleRepo "PersonalBlog/internal/repositories/article"
	"context"
	"github.com/google/uuid"
	"time"
)

type ArticleUseCase interface {
	GetAllArticles(ctx context.Context) ([]domain.Article, error)
	GetArticleByID(ctx context.Context, id string) (*domain.Article, error)
	CreateArticle(ctx context.Context, article *domain.Article) error
	UpdateArticle(ctx context.Context, id string, article *domain.Article) error
	DeleteArticle(ctx context.Context, id string) error
}

type articleUseCase struct {
	repo articleRepo.ArticleRepository
}

func NewArticleUseCase(repo articleRepo.ArticleRepository) ArticleUseCase {
	return &articleUseCase{repo: repo}
}

func (uc *articleUseCase) GetAllArticles(ctx context.Context) ([]domain.Article, error) {
	return uc.repo.GetAll(ctx)
}

func (uc *articleUseCase) GetArticleByID(ctx context.Context, id string) (*domain.Article, error) {
	return uc.repo.GetByID(ctx, id)
}

func (uc *articleUseCase) CreateArticle(ctx context.Context, article *domain.Article) error {
	article.ID = uuid.New().String()
	article.CreatedAt = time.Now()
	article.UpdatedAt = time.Now()

	return uc.repo.Create(ctx, article)
}

func (uc *articleUseCase) UpdateArticle(ctx context.Context, id string, article *domain.Article) error {
	article.UpdatedAt = time.Now()

	return uc.repo.Update(ctx, id, article)
}

func (uc *articleUseCase) DeleteArticle(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
