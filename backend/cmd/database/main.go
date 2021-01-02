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

const databaseVersion = 3

// Init establishes a database connection, applies migrations, and returns a *sql.DB instance
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
		"file://cmd/database/migrations", "postgres", driver)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	migrationVersion, dirty, err := m.Version()
	if dirty {
		logger.Fatal("the current database schema is reported as being dirty. A manual resolution is needed")
	}
	if err != nil {
		return nil, err
	}
	logger.Infof("migration - migrationVersion: %v dirty: %v", migrationVersion, dirty)
	if databaseVersion != migrationVersion {
		err = nil
		if databaseVersion > migrationVersion {
			err = m.Up()
		} else {
			err = m.Down()
		}

		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	logger.Info("database connection/migration successful")
	return db, nil
}
