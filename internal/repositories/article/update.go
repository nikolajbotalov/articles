package article

import (
	"PersonalBlog/internal/domain"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *articleRepository) Update(ctx context.Context, id string, article *domain.Article) error {
	r.logger.Info("Updating article", zap.String("id", id))

	// Создаем SQL-запрос
	query, args, err := sq.Update("articles").
		Set("author", article.Author).
		Set("title", article.Title).
		Set("content", article.Content).
		Set("updated_at", article.UpdatedAt).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error("Failed to build SQL query", zap.Error(err))
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к БД
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to update article", zap.Error(err))
		return fmt.Errorf("failed to update article: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	if result.RowsAffected() == 0 {
		r.logger.Error("Article not found to update", zap.String("id", id))
		return fmt.Errorf("article not found")
	}

	r.logger.Info("Successfully updated article", zap.String("id", id))
	return nil
}
