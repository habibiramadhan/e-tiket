package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	
	"ticket-system/internal/delivery/http/middleware"
	"ticket-system/internal/delivery/http/routes"
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

	errorMiddleware := middleware.NewErrorMiddleware()
	
	app := fiber.New(fiber.Config{
		AppName:       "Ticket System API",
		CaseSensitive: true,
		StrictRouting: true,
		ServerHeader:  "Fiber",
		ErrorHandler:  errorMiddleware.ErrorHandler(),
	})
	
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:8080", 
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowCredentials: true,
	}))
	
	routes.SetupRoutes(app, db, cfg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	go func() {
		<-c
		log.Println("Shutting down server...")
		_ = app.Shutdown()
	}()

	serverAddr := fmt.Sprintf(":%s", cfg.ServerPort)
	log.Printf("Server berjalan di http://localhost%s", serverAddr)
	
	if err := app.Listen(serverAddr); err != nil {
		log.Fatalf("Gagal memulai server: %v", err)
	}
}