package migrations

import (
	"fmt"

	"github.com/Helltale/beer-mania/backend/internal/entity"
	"gorm.io/gorm"
)

// RunMigrations runs all database migrations using GORM AutoMigrate
func RunMigrations(db *gorm.DB) error {
	// Enable UUID extension for PostgreSQL
	if err := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error; err != nil {
		return fmt.Errorf("failed to create uuid extension: %w", err)
	}

	// AutoMigrate creates tables and indexes based on entity definitions
	if err := db.AutoMigrate(
		&entity.Image{},
		&entity.ProcessingTask{},
	); err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	// Create foreign key constraint for processing_tasks.image_id -> images.id
	// GORM doesn't automatically create foreign keys with AutoMigrate,
	// so we need to create it manually if it doesn't exist
	if err := db.Exec(`
		DO $$
		BEGIN
			IF NOT EXISTS (
				SELECT 1 FROM pg_constraint 
				WHERE conname = 'fk_processing_tasks_image_id'
			) THEN
				ALTER TABLE processing_tasks 
				ADD CONSTRAINT fk_processing_tasks_image_id 
				FOREIGN KEY (image_id) 
				REFERENCES images(id) 
				ON DELETE CASCADE;
			END IF;
		END $$;
	`).Error; err != nil {
		return fmt.Errorf("failed to create foreign key constraint: %w", err)
	}

	return nil
}

// RollbackMigrations drops all tables (use with caution!)
func RollbackMigrations(db *gorm.DB) error {
	if err := db.Migrator().DropTable(
		&entity.ProcessingTask{},
		&entity.Image{},
	); err != nil {
		return fmt.Errorf("failed to rollback migrations: %w", err)
	}
	return nil
}
