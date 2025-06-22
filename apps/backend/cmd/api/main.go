package main

import (
	"log"
	"os"

	"github.com/davidbyttow/govips/v2/vips"
	"github.com/joho/godotenv"
	database "github.com/okzmo/kyob/db"
	"github.com/okzmo/kyob/internal/api/actors"
	"github.com/okzmo/kyob/internal/api/router"
	"github.com/okzmo/kyob/internal/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db := database.Setup()
	defer db.Close()

	utils.SetupSnowflake()

	if len(os.Args) > 1 && os.Args[1] == "seed" {
		database.RunSeeder()
		return
	}

	vips.Startup(nil)
	defer vips.Shutdown()

	actors.SetupServersEngine()
	actors.SetupUsersEngine()
	router.Setup()
}
