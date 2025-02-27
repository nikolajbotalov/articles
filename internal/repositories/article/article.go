package repositories

import (
	"PersonalBlog/internal/domain"
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArticleRepository interface {
	GetAll(ctx context.Context) ([]domain.Article, error)
	GetByID(ctx context.Context, id string) (*domain.Article, error)
	Create(ctx context.Context, article *domain.Article) error
	Update(ctx context.Context, id string, article *domain.Article) error
	Delete(ctx context.Context, id string) error
}

type articleRepository struct {
	db *pgxpool.Pool
}

func NewArticleRepository(db *pgxpool.Pool) ArticleRepository {
	return &articleRepository{db: db}
}

func (r *articleRepository) GetAll(ctx context.Context) ([]domain.Article, error) {
	// Создаем запрос с помощью squirrel
	query, args, err := sq.Select("id", "author", "title", "content", "created_at", "updated_at").
		From("articles").PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к базе данных
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		articles = append(articles, article)
	}

	// Проверяем, были ли ошибки при итерации по строкам
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return articles, nil
}

func (r *articleRepository) GetByID(ctx context.Context, id string) (*domain.Article, error) {
	// Создаем SQL-запрос
	query, args, err := sq.Select("id", "author", "title", "content", "created_at", "updated_at").
		From("articles").Where(sq.Eq{"id": id}).PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
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
			return nil, fmt.Errorf("article not found: %w", err)
		}
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &article, nil
}

func (r *articleRepository) Create(ctx context.Context, article *domain.Article) error {
	// Создаем SQL-запрос
	query, args, err := sq.Insert("articles").
		Columns("id", "author", "title", "content", "created_at", "updated_at").
		Values(article.ID, article.Author, article.Title, article.Content, article.CreatedAt, article.UpdatedAt).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к базе данных
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to insert article into database: %w", err)
	}

	return nil
}

func (r *articleRepository) Update(ctx context.Context, id string, article *domain.Article) error {
	// Создаем SQL-запрос
	query, args, err := sq.Update("articles").
		Set("author", article.Author).
		Set("title", article.Title).
		Set("content", article.Content).
		Set("updated_at", article.UpdatedAt).
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к БД
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update article: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	if result.RowsAffected() == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}

func (r *articleRepository) Delete(ctx context.Context, id string) error {
	// Создаем SQL-запрос
	query, args, err := sq.Delete("articles").Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar).ToSql()
	if err != nil {
		return fmt.Errorf("failed to build SQL query: %w", err)
	}

	// Выполняем запрос к БД
	result, err := r.db.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete article: %w", err)
	}

	// Проверяем, была ли обновлена хотя бы одна строка
	if result.RowsAffected() == 0 {
		return fmt.Errorf("article not found")
	}

	return nil
}
