package db

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// Init connects to PostgreSQL and runs auto-migration.
func Init() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(fmt.Sprintf("failed to connect to database: %v", err))
	}

	// Auto-migrate schema
	if err := DB.AutoMigrate(&Match{}, &MatchResult{}); err != nil {
		log.Fatal(fmt.Sprintf("failed to migrate database: %v", err))
	}

	log.Println("Database connected and migrated successfully")
}
