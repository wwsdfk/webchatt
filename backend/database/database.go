package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func InitDB() error {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	DB, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		return fmt.Errorf("ошибка подключения к базе данных: %w", err)
	}

	if err = DB.Ping(ctx); err != nil {
		DB.Close()
		return fmt.Errorf("ошибка пинга базы данных: %w", err)
	}

	log.Println("Успешное подключение к базе данных")

	_, err = DB.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS messages (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			content TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		DB.Close()
		return fmt.Errorf("ошибка создания таблицы: %w", err)
	}

	return nil
}
