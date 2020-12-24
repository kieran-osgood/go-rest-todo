package main

import (
	"fmt"
	"log"

	"github.com/joho/godotenv"
	"github.com/kieran-osgood/go-rest-todo/api"
	"github.com/kieran-osgood/go-rest-todo/config"
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
	config := config.New()
	pgsql := database.Database{
		Host:   config.Database.Host,
		Port:   config.Database.Port,
		User:   config.Database.User,
		Pass:   config.Database.Pass,
		DbName: config.Database.DbName,
	}

	err := pgsql.Init()
	if err != nil {
		fmt.Println("pgsql.Init")
		panic(err)
	}

	err = api.Init()
	if err != nil {
		fmt.Println("api.Init")
		panic(err)
	}

}
