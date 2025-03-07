//pkg/database/migrate.go

package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

func MigrateDatabase(db *sql.DB, migrationPath string) error {
	log.Printf("Menjalankan migrasi dari: %s", migrationPath)

	content, err := ioutil.ReadFile(migrationPath)
	if err != nil {
		return fmt.Errorf("gagal membaca file migrasi: %w", err)
	}

	sqlContent := string(content)

	re := regexp.MustCompile(`(?i)(CREATE OR REPLACE FUNCTION.*?LANGUAGE\s+[a-z]+;|CREATE\s+TRIGGER.*?END;|CREATE.*?;)`)
	matches := re.FindAllString(sqlContent, -1)

	for i, statement := range matches {
		if statement == "" {
			continue
		}

		if _, err := db.Exec(statement); err != nil {
			return fmt.Errorf("gagal mengeksekusi query ke-%d: %w", i+1, err)
		}
		log.Printf("Query ke-%d berhasil dieksekusi", i+1)
	}

	log.Println("Migrasi database berhasil dijalankan")
	return nil
}