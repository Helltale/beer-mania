package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Helltale/beer-mania/backend/migrations"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	var (
		dsn      = flag.String("dsn", "", "Database connection string (required)")
		rollback = flag.Bool("rollback", false, "Rollback migrations (drops all tables)")
	)
	flag.Parse()

	if *dsn == "" {
		// Try to get DSN from environment variables
		host := getEnv("POSTGRES_HOST", "localhost")
		port := getEnv("POSTGRES_PORT", "5432")
		user := getEnv("POSTGRES_USER", "beermania_user")
		password := getEnv("POSTGRES_PASSWORD", "beermania_password")
		dbname := getEnv("POSTGRES_DB", "beermania_db")
		sslmode := getEnv("POSTGRES_SSLMODE", "disable")

		*dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			host, port, user, password, dbname, sslmode)
	}

	// Connect to database
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	if *rollback {
		log.Println("Rolling back migrations...")
		if err := migrations.RollbackMigrations(db); err != nil {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
		return
	}

	log.Println("Running migrations...")
	if err := migrations.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
