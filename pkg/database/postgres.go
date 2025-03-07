//pkg/database/postgres.go

package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresConnection(config PostgresConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	log.Println("Berhasil terhubung ke database PostgreSQL")
	return db, nil
}

func ClosePostgresConnection(db *sql.DB) {
	if db != nil {
		if err := db.Close(); err != nil {
			log.Printf("Error saat menutup koneksi database: %v", err)
			return
		}
		log.Println("Koneksi database PostgreSQL berhasil ditutup")
	}
}