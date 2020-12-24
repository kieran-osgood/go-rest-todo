package database

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	_ "github.com/lib/pq"
	"go.uber.org/zap"
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
func (d *Database) Init(logger *zap.SugaredLogger) (*sql.DB, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Pass, d.DbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		logger.Error(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations", "postgres", driver)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	version, dirty, err := m.Version()
	logger.Infof("migration - version: %v dirty: %v", version, dirty)

	m.Steps(3)

	logger.Info("database connection/migration successful")
	return db, nil
}
