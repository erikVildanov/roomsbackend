package storage

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"roomsbackend/services/common/logger"
)

// Config описывает конфигурацию подключения к базе данных PostgreSQL.
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// NewPostgres подключается к PostgreSQL с заданной конфигурацией.
func NewPostgres(cfg Config) *sql.DB {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Database,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		logger.Errorf("❌ Ошибка подключения к PostgreSQL: %v", err)
	}

	if err := db.Ping(); err != nil {
		logger.Errorf("❌ PostgreSQL не отвечает: %v", err)
	}

	logger.Infof("✅ Подключение к PostgreSQL успешно")
	return db
}
