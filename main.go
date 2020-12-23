package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kieran-osgood/go-rest-todo/api"
	"github.com/kieran-osgood/go-rest-todo/database"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	database.Init()
	api.Init()

}
