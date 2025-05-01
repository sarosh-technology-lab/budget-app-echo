package main

import (
	"budget-backend/internal/database/seeders"
)

func main() {
	seeders.RunAllSeeders()
}