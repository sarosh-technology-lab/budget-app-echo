package main

import (
	"budget-backend/common"
	"budget-backend/internal/models"
	"log"
)

func main() {
	db, err := common.Mysql()
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&models.User{})
	if err != nil {
		panic(err)
	}
	log.Println("Migration completed")
}