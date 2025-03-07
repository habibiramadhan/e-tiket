//cmd/api/main.go

package main

import (
	"fmt"
	"log"
	"net/http"

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

	if err := db.Ping(); err != nil {
		log.Fatalf("Gagal melakukan ping ke database: %v", err)
	}
	
	log.Println("Koneksi ke database berhasil")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Aplikasi Event Management API berjalan")
	})

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server berjalan di http://localhost%s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}