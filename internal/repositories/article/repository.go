package article

import (
	"PersonalBlog/internal/domain"
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository interface {
	GetAll(ctx context.Context) ([]domain.Article, error)
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	Create(ctx context.Context, article *domain.Article) error
	Update(ctx context.Context, id string, article *domain.Article) error
	Delete(ctx context.Context, id string) error
}

type articleRepository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

func NewArticleRepository(db *pgxpool.Pool, logger *zap.Logger) Repository {
	return &articleRepository{
		db:     db,
		logger: logger,
	}
}
