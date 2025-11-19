package postgres

import (
	"context"
	dbsql "database/sql"
	"fmt"
	"golang_template/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
	pgorm "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Database interface {
	Close() error
	DB() *dbsql.DB
	Gorm() *gorm.DB
}

type database struct {
	pool     *pgxpool.Pool
	database *dbsql.DB
	gormDB   *gorm.DB
}

func NewDatabase(ctx context.Context, cfg *config.DatabaseConfig, logger *zap.Logger) (Database, error) {
	if cfg == nil {
		return nil, fmt.Errorf("database config cannot be nil")
	}

	if cfg.MaxConns < cfg.MinConns {
		return nil, fmt.Errorf("maxConns must be greater than or equal to minConns")
	}

	poolConfig, err := pgxpool.ParseConfig(config.GetDSN(cfg))
	if err != nil {
		return nil, err
	}

	// Configure pool
	poolConfig.MaxConns = int32(cfg.MaxConns)
	poolConfig.MinConns = int32(cfg.MinConns)

	timeoutCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(timeoutCtx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("creating pool: %w", err)
	}

	// Test connection
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging database: %w", err)
	}

	sqlDB := stdlib.OpenDB(*pool.Config().ConnConfig)

	gormDB, err := gorm.Open(pgorm.New(pgorm.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		pool.Close()
		sqlDB.Close()
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	gormDB.AutoMigrate()

	return &database{
		pool:     pool,
		gormDB:   gormDB,
		database: sqlDB,
	}, nil
}

func (db *database) Close() error {
	var errs []error

	db.pool.Close()
	if err := db.database.Close(); err != nil {
		errs = append(errs, fmt.Errorf("closing database: %w", err))
	}

	if len(errs) > 0 {
		return fmt.Errorf("closing database: %v", errs)
	}
	return nil
}

func (db *database) Gorm() *gorm.DB {
	return db.gormDB
}

func (db *database) DB() *dbsql.DB {
	return db.database
}
