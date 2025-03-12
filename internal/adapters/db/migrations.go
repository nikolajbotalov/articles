package db

import (
	"PersonalBlog/internal/config"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/golang-migrate/migrate/v4"
	"go.uber.org/zap"
	"time"
)

func RunMigrations(cfg config.PostgreSQL, logger *zap.Logger) error {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	logger.Debug("Attempting to connect to DB", zap.String("dsn", dsn))

	var m *migrate.Migrate
	err := retry.Do(
		func() error {
			var err error
			m, err = migrate.New("file://migrations", dsn)
			if err != nil {
				logger.Warn("DB not ready, retrying...", zap.Error(err))
				return err
			}
			return nil
		},
		retry.Attempts(10),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migration: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}
	logger.Info("Migrations applied successfully")
	return nil
}
