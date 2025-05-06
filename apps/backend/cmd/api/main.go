package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/router"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := db.Setup()
	defer db.Close()

	router.Setup()
}
