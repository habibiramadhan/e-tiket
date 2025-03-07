//pkg/database/simple_migrate.go

package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func SimpleMigrateDatabase(db *sql.DB, migrationPath string) error {
	log.Printf("Menjalankan migrasi dari: %s", migrationPath)

	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("gagal membaca file migrasi: %w", err)
	}

	sqlContent := string(content)

	if _, err := db.Exec(sqlContent); err != nil {
		return fmt.Errorf("gagal mengeksekusi SQL: %w", err)
	}

	log.Println("Migrasi database berhasil dijalankan")
	return nil
}