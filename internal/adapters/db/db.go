package db

import (
	"PersonalBlog/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
	"strconv"
)

type PostgreSQLDB struct {
	pool   *pgxpool.Pool
	logger *zap.Logger
}

func NewPostgreSQLDB(cfgPSQL config.PostgreSQL, logger *zap.Logger) (*pgxpool.Pool, error) {
	logger.Info("Connecting to PostgreSQL")

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

	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		logger.Error("Cannot connect to PostgreSQL", zap.Error(err))
		return nil, err
	}

	// Проверяем, что соединение с базой данных работает
	if err := pool.Ping(context.Background()); err != nil {
		logger.Error("Cannot ping PostgreSQL", zap.Error(err))
		return nil, err
	}

	logger.Info("Connected to PostgreSQL")
	return pool, nil
}

func (db *PostgreSQLDB) Close() error {
	if db.pool != nil {
		db.pool.Close()
	}
	return nil
}
