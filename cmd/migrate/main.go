//cmd/migrate/main.go

package main

import (
	"log"
	"path/filepath"

	"ticket-system/pkg/config"
    "ticket-system/pkg/database"
)

func main() {
	cfg := config.LoadConfig()

	dbConfig := database.PostgresConfig{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  cfg.DBSSLMode,
	}

	db, err := database.NewPostgresConnection(dbConfig)
	if err != nil {
		log.Fatalf("Gagal menginisialisasi database: %v", err)
	}
	defer database.ClosePostgresConnection(db)

	migrationPath := filepath.Join("migrations", "schema.sql")

	if err := database.SimpleMigrateDatabase(db, migrationPath); err != nil {
		log.Printf("Migrasi sederhana gagal: %v", err)
		log.Println("Mencoba dengan metode alternatif...")
		
		if err := database.MigrateDatabase(db, migrationPath); err != nil {
			log.Fatalf("Gagal menjalankan migrasi: %v", err)
		}
	}

	log.Println("Migrasi selesai")
}