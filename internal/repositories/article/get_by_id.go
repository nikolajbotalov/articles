package article

import (
	"PersonalBlog/internal/domain"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

func (r *articleRepository) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	r.logger.Info("Fetching article by ID", zap.String("id", id))

	// Создаем SQL-запрос
	query, args, err := sq.Select("id", "author", "title", "content", "created_at", "updated_at").
		From("articles").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error("Failed to build query", zap.Error(err))
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к БД
	row := r.db.QueryRow(ctx, query, args...)

	// Сканируем результат в структуру Article
	var article domain.Article
	if err := row.Scan(
		&article.ID,
		&article.Author,
		&article.Title,
		&article.Content,
		&article.CreatedAt,
		&article.UpdatedAt,
	); err != nil {
		if err == pgx.ErrNoRows {
			r.logger.Error("Article not found", zap.String("id", id))
			return nil, fmt.Errorf("article not found: %w", err)
		}
		r.logger.Error("Failed to scan row", zap.Error(err))
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	r.logger.Info("Successfully fetched article by ID", zap.String("id", id))
	return &article, nil
}
