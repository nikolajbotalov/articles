package app

import (
	"PersonalBlog/internal/adapters/db"
	"PersonalBlog/internal/config"
	"PersonalBlog/internal/logger"
	repositories "PersonalBlog/internal/repositories/article"
	usecases "PersonalBlog/internal/usecases/article"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	Config *config.Config
	DB     *db.PostgreSQLDB
	Server *Server
}

// NewApp инициализирует и возвращает новое приложение
func NewApp() (*App, error) {
	// инициализация логера
	zapLogger, err := logger.NewLogger()
	if err != nil {
		return nil, err
	}

	// инициализация конфига
	cfg := config.GetConfig()
	zapLogger.Debug("Config loaded",
		zap.String("db_host", cfg.PostgreSQL.Host),
		zap.String("db_port", cfg.PostgreSQL.Port),
		zap.String("db_username", cfg.PostgreSQL.Username),
		zap.String("db_database", cfg.PostgreSQL.Database))

	// запуск миграций
	if err := db.RunMigrations(cfg.PostgreSQL, zapLogger); err != nil {
		zapLogger.Error("Failed to run migrations", zap.Error(err))
		return nil, err
	}

	// иниализация БД
	dbInstance, err := db.NewPostgreSQLDB(cfg.PostgreSQL, zapLogger)
	if err != nil {
		zapLogger.Error("failed to initialize db", zap.Error(err))
		return nil, err
	}

	articleRepo := repositories.NewArticleRepository(dbInstance.Pool(), zapLogger)
	articleUseCase := usecases.NewArticleUseCase(articleRepo, zapLogger)

	// инициализация сервера
	server := NewServer(cfg, articleUseCase, zapLogger)

	return &App{
		Logger: zapLogger,
		Config: cfg,
		DB:     dbInstance,
		Server: server,
	}, nil
}

// Close освобождает ресурсы приложения
func (a *App) Close() {
	a.Logger.Info("Closing application")
	if err := a.DB.Close(); err != nil {
		a.Logger.Error("Failed to close DB", zap.Error(err))
	}
}
