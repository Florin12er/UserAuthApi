package database

import (
	"log"
	"UserAuth/internal/models"
)

func SyncDatabase() {
	// Check if the users table exists
	if DB.Migrator().HasTable(&models.User{}) {
		log.Println("Users table already exists. Migrating schema.")
		// AutoMigrate will only add missing columns and indexes, it won't delete/change existing columns
		if err := DB.AutoMigrate(&models.User{}); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	} else {
		// If the table doesn't exist, create it
		log.Println("Creating users table.")
		if err := DB.AutoMigrate(&models.User{}); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	}
}
