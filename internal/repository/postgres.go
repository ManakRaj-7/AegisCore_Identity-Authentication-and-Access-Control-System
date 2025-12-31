package repository

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"go.uber.org/zap"

	"github.com/randhir/aegis-core/internal/config"
	"github.com/randhir/aegis-core/internal/logger"
)

var DB *sql.DB

func ConnectPostgres() error {
	cfg := config.AppConfig.Database

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.Password,
		cfg.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %w", err)
	}

	DB = db

	logger.Info("PostgreSQL connected successfully",
		zap.String("host", cfg.Host),
		zap.String("port", cfg.Port),
		zap.String("database", cfg.Name),
	)

	return nil
}

func ClosePostgres() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}

