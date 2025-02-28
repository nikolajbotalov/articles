package article

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *articleRepository) Delete(ctx context.Context, id string) error {
	r.logger.Info("Deleting article", zap.String("id", id))

	// Создаем SQL-запрос
	query, args, err := sq.Delete("articles").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error("Failed to build query", zap.Error(err))
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к БД
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to delete article", zap.Error(err))
		return fmt.Errorf("failed to delete article: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	if result.RowsAffected() == 0 {
		r.logger.Error("Article not found for deletion", zap.String("id", id))
		return fmt.Errorf("article not found")
	}

	r.logger.Info("Successfully deleted article", zap.String("id", id))
	return nil
}
