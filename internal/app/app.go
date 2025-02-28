package app

import (
	"PersonalBlog/internal/adapters/db"
	"PersonalBlog/internal/config"
	"PersonalBlog/internal/logger"
	repositories "PersonalBlog/internal/repositories/article"
	usecases "PersonalBlog/internal/usecases/article"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type App struct {
	Logger *zap.Logger
	Config *config.Config
	DB     *pgxpool.Pool
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
	zapLogger.Debug("config loaded", zap.Any("config", cfg))

	// иниализация БД
	dbPool, err := db.NewPostgreSQLDB(cfg.PostgreSQL)
	if err != nil {
		zapLogger.Error("failed to initialize db", zap.Error(err))
		return nil, err
	}

	articleRepo := repositories.NewArticleRepository(dbPool, zapLogger)
	articleUseCase := usecases.NewArticleUseCase(articleRepo)

	// инициализация сервера
	server := NewServer(cfg, articleUseCase, zapLogger)

	return &App{
		Logger: zapLogger,
		Config: cfg,
		DB:     dbPool,
		Server: server,
	}, nil
}

// Close освобождает ресурсы приложения
func (a *App) Close() {
	a.Logger.Info("Closing application")
	a.DB.Close()
}
