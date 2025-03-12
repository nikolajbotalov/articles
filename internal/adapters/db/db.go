package db

import (
	"PersonalBlog/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type PostgreSQLDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgreSQLDB(cfgPSQL config.PostgreSQL, logger *zap.Logger) (*PostgreSQLDB, error) {
	logger.Info("Connecting to PostgreSQL", zap.String("host", cfgPSQL.Host), zap.String("port", cfgPSQL.Port))

	dbPort, err := strconv.ParseUint(cfgPSQL.Port, 10, 16)
	if err != nil {
		logger.Error("Cannot convert DB port to int", zap.Error(err))
		return nil, err
	}

	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		cfgPSQL.Host,
		dbPort,
		cfgPSQL.Username,
		cfgPSQL.Password,
		cfgPSQL.Database)

	const maxAttempts = 10
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		pool, err := pgxpool.New(context.Background(), connStr)
		if err != nil {
			logger.Warn("Failed to create connection pool", zap.Int("attempt", attempt), zap.Error(err))
			if attempt == maxAttempts {
				logger.Error("Cannot connect to PostgreSQL after retries", zap.Error(err))
				return nil, err
			}
			time.Sleep(2 * time.Second)
			continue
		}

		// Проверяем, что соединение с базой данных работает
		if err := pool.Ping(context.Background()); err != nil {
			logger.Warn("Failed to ping PostgreSQL", zap.Int("attempt", attempt), zap.Error(err))
			pool.Close()
			if attempt == maxAttempts {
				logger.Error("Cannot ping PostgreSQL after retries", zap.Error(err))
				return nil, err
			}
			time.Sleep(2 * time.Second)
			continue
		}

		logger.Info("Connected to PostgreSQL")
		return &PostgreSQLDB{
			pool:   pool,
			logger: logger,
		}, nil
	}

	return nil, fmt.Errorf("failed to connect to PostgreSQL after %d attempts", maxAttempts)
}

func (db *PostgreSQLDB) Close() error {
	if db.pool != nil {
		db.pool.Close()
	}
	return nil
}

func (db *PostgreSQLDB) Pool() *pgxpool.Pool {
	return db.pool
}
