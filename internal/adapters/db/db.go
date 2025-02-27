package db

import (
	"PersonalBlog/internal/config"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"strconv"
)

type PostgreSQLDB struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLDB(cfgPSQL config.PostgreSQL) (*pgxpool.Pool, error) {
	dbPort, err := strconv.ParseUint(cfgPSQL.Port, 10, 16)
	if err != nil {
		slog.Error("Cannot convert DB port to int: %s", err)
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
		slog.Error("Cannot create DB pool: %v", err)
		return nil, err
	}

	// Проверяем, что соединение с базой данных работает
	if err := pool.Ping(context.Background()); err != nil {
		slog.Error("Can't ping DB: %v", err)
		return nil, err
	}

	return pool, nil
}

func (db *PostgreSQLDB) Close() error {
	if db.pool != nil {
		db.pool.Close()
	}
	return nil
}
