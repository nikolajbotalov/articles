package article

import (
	"PersonalBlog/internal/domain"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

func (r *articleRepository) GetAll(ctx context.Context) ([]domain.Article, error) {
	r.logger.Info("Fetching all articles")

	// Создаем запрос с помощью squirrel
	query, args, err := sq.Select("id", "author", "title", "content", "created_at", "updated_at").
		From("articles").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		r.logger.Error("Failed to build SQL query", zap.Error(err))
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к базе данных
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute SQL query", zap.Error(err))
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}
	defer rows.Close()

	// Сканируем результаты в структуру Article
	var articles []domain.Article
	for rows.Next() {
		var article domain.Article
		if err := rows.Scan(
			&article.ID,
			&article.Author,
			&article.Title,
			&article.Content,
			&article.CreatedAt,
			&article.UpdatedAt,
		); err != nil {
			r.logger.Error("Failed to scan row", zap.Error(err))
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		articles = append(articles, article)
	}

	// Проверяем, были ли ошибки при итерации по строкам
	if err := rows.Err(); err != nil {
		r.logger.Error("error while iterating over rows", zap.Error(err))
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	r.logger.Info("Successfully fetched all articles", zap.Int("count", len(articles)))
	return articles, nil
}
