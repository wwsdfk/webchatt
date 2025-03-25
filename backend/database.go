package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func initDB() {
	dbURL := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err = pgxpool.New(ctx, dbURL)
	if err != nil {
		log.Fatalf("Ошибка подключения к базе данных: %v", err)
	}

	if err = db.Ping(ctx); err != nil {
		log.Fatalf("Ошибка пинга базы данных: %v", err)
	}

	fmt.Println("Успешное подключение к базе данных")
}
