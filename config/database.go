package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	DB      *gorm.DB
	PgxConn *pgx.Conn
)

func ConnectDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Build connection string
	databaseUrl := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s&search_path=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
		os.Getenv("DB_SEARCH_PATH"),
	)

	// Test connection with pgx
	pgxConn, err := pgx.Connect(context.Background(), databaseUrl)
	if err != nil {
		log.Fatal("Failed to connect with pgx: ", err)
	}

	// Test connection
	var version string
	err = pgxConn.QueryRow(context.Background(), "SELECT version()").Scan(&version)
	if err != nil {
		log.Fatal("Failed to query version: ", err)
	}
	log.Println("PostgreSQL version:", version)

	PgxConn = pgxConn

	// Connect with GORM using standard DSN
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		PrepareStmt: false,
	})
	if err != nil {
		log.Fatal("Failed to connect database with GORM: ", err)
	}

	DB = db
	log.Println("Database connected successfully")
}

func CloseDB() {
	if PgxConn != nil {
		PgxConn.Close(context.Background())
	}
}
