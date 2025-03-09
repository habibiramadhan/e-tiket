//pkg/config/config.go

package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost      string
	DBPort      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBSSLMode   string
	ServerPort  string
	JWTSecret   string
	TokenExpiry string
	AppEnv      string
	
	// SMTP Settings
	SMTPHost     string
	SMTPPort     string
	SMTPUsername string
	SMTPPassword string
	SMTPFromName string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan")
	}

	config := &Config{
		DBHost:      getEnv("DB_HOST", ""),
		DBPort:      getEnv("DB_PORT", ""),
		DBUser:      getEnv("DB_USER", ""),
		DBPassword:  getEnv("DB_PASSWORD", ""),
		DBName:      getEnv("DB_NAME", ""),
		DBSSLMode:   getEnv("DB_SSLMODE", ""),
		ServerPort:  getEnv("SERVER_PORT", ""),
		JWTSecret:   getEnv("JWT_SECRET", "rahasia_aku_kamu_dan_jwt"),
		TokenExpiry: getEnv("TOKEN_EXPIRY", "24"),
		AppEnv:      getEnv("APP_ENV", "development"),
		
		// SMTP Settings
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnv("SMTP_PORT", ""),
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFromName: getEnv("SMTP_FROM_NAME", "Sistem Tiket Event"),
	}

	return config
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}