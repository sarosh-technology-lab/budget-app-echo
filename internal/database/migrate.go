package main

import (
	"budget-backend/common"
	"budget-backend/internal/models"
	"flag"
	"log"
	"strings"
	"gorm.io/gorm"
)

func migrateUp(db *gorm.DB) error {
 // Create enum type first
    err := db.Exec("CREATE TYPE gender_enum AS ENUM ('M', 'F', 'O')").Error
    if err != nil {
        // Ignore error if enum already exists
        if !strings.Contains(err.Error(), "already exists") {
            return err
        }
    }

	// Up migration: Create tables
	err = db.AutoMigrate(&models.User{}, &models.AppToken{}, &models.Category{}, &models.Role{}, &models.Permission{}, &models.RolePermission{})
	if err != nil {
		return err
	}
	log.Println("Up migration completed (tables created)")
	return nil
}

func migrateDown(db *gorm.DB) error {
	// Down migration: Drop tables
	err := db.Migrator().DropTable(&models.User{}, &models.AppToken{}, &models.Category{}, &models.Role{}, &models.Permission{}, &models.RolePermission{})
	if err != nil {
		return err
	}

	// Drop the enum type
    err = db.Exec("DROP TYPE IF EXISTS gender_enum").Error
    if err != nil {
        return err
    }
	log.Println("Down migration completed (tables dropped)")
	return nil
}

func main() {
	// Define a flag to choose migration direction
	direction := flag.String("direction", "up", "Migration direction: 'up' or 'down'")
	flag.Parse()

	// Connect to the database
	db, err := common.Sql()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Run the appropriate migration based on the flag
	switch *direction {
	case "up":
		if err := migrateUp(db); err != nil {
			log.Fatalf("Up migration failed: %v", err)
		}
	case "down":
		if err := migrateDown(db); err != nil {
			log.Fatalf("Down migration failed: %v", err)
		}
	default:
		log.Fatal("Invalid direction. Use '-direction=up' or '-direction=down'")
	}
}