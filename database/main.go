package database

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/kieran-osgood/go-rest-todo/config"
	_ "github.com/lib/pq"
)

func Init() {
	config := config.New()

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.Database.Host, config.Database.Port, config.Database.User, config.Database.Pass, config.Database.DbName)
	log.Printf(psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
}
