package main

import (
	"log"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/joho/godotenv"
	"github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	"github.com/okzmo/kyob/internal/api/router"
	"github.com/okzmo/kyob/internal/utils"
)

func main() {
	vips.Startup(nil)
	defer vips.Shutdown()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := db.Setup()
	defer db.Close()

	utils.SetupSnowflake()

	actors.SetupServersEngine()
	actors.SetupUsersEngine()
	router.Setup()
}
