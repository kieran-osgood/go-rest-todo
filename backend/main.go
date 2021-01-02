package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kieran-osgood/go-rest-todo/cmd/api"
	"github.com/kieran-osgood/go-rest-todo/cmd/config"
	"github.com/kieran-osgood/go-rest-todo/cmd/database"
	errorHandler "github.com/kieran-osgood/go-rest-todo/cmd/error"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Panic("./main.go - Can't find .env file")
	}
}

func initLogger() (*zap.SugaredLogger, error) {
	loggerConfig := zap.NewDevelopmentConfig()
	loggerConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	sugar := logger.Sugar()

	return sugar, nil
}

func main() {
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	defer errorHandler.CleanUpAndHandleError(logger.Sync, logger)

	c := config.New()
	pgsql := database.Database{
		Host:   c.Database.Host,
		Port:   c.Database.Port,
		User:   c.Database.User,
		Pass:   c.Database.Pass,
		DbName: c.Database.DbName,
	}

	db, err := pgsql.Init(logger)
	if err != nil {
		logger.Panicf("pgsql init: %v", err)
	}
	defer errorHandler.CleanUpAndHandleError(db.Close, logger)

	err = api.Init(logger, db)
	if err != nil {
		logger.Panicf("api init: %v", err)
	}
}
