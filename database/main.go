package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	_ "github.com/lib/pq"
)

// Database constructor
type Database struct {
	Host   string
	Port   int
	User   string
	Pass   string
	DbName string
}

const version = 1

// Init function for database
func (d *Database) Init() error {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Pass, d.DbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return err
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		return err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres", driver)

	if err != nil {
		return err
	}

	// version, dirty, err := m.Version()

	// zap.L().Infof("Migration - Version: %v Dirty: %v \n", version, dirty)

	m.Steps(3)

	fmt.Println("Successfully connected!")
	return nil
}
