package database

import (
	"fmt"
	"time"

	"anarchy-core/internal/util" // Импорт вашего пакета util

	"github.com/jmoiron/sqlx" // Для удобной работы с SQL
	_ "github.com/lib/pq"     // PostgreSQL драйвер
)

// InitPostgresDB initializes and returns a new PostgreSQL database connection.
func InitPostgresDB(databaseURL string, logger *util.Logger) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Set connection pool settings
	db.SetMaxOpenConns(25)                 // Максимальное количество открытых соединений
	db.SetMaxIdleConns(10)                 // Максимальное количество неактивных соединений
	db.SetConnMaxLifetime(5 * time.Minute) // Максимальное время жизни соединения

	// Ping the database to verify the connection
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Info("Successfully connected to PostgreSQL database!")
	return db, nil
}
