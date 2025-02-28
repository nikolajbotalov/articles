package article

import (
	"PersonalBlog/internal/domain"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *articleRepository) Create(ctx context.Context, article *domain.Article) error {
	r.logger.Info("Create article", zap.String("title", article.Title))

	// Создаем SQL-запрос
	query, args, err := sq.Insert("articles").
		Columns("id", "author", "title", "content", "created_at", "updated_at").
		Values(article.ID, article.Author, article.Title, article.Content, article.CreatedAt, article.UpdatedAt).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error("Failed to build query", zap.Error(err))
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к базе данных
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to insert article into database", zap.Error(err))
		return fmt.Errorf("failed to insert article into database: %w", err)
	}

	r.logger.Info("Successfully created article", zap.String("id", article.ID))
	return nil
}
