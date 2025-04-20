package common

import (
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Sql() (*gorm.DB, error) {

	err := godotenv.Load()

	if err != nil {
		panic("Error Loading .env file")
	}

	host := os.Getenv("DB_HOST")
	database := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", 
        host, username, password, database, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{TranslateError: true})
	if err != nil {
		return nil, err
	}
	log.Default().Println("Database Connection Successful")

	return db, nil
}
