package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/kieran-osgood/go-rest-todo/api"
	"github.com/kieran-osgood/go-rest-todo/config"
	"github.com/kieran-osgood/go-rest-todo/database"
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
	config := zap.NewDevelopmentConfig()
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar, nil
}

func main() {
	logger, err := initLogger()
	if err != nil {
		log.Fatalf("can't initialize zap logger: %v", err)
	}

	config := config.New()

	pgsql := database.Database{
		Host:   config.Database.Host,
		Port:   config.Database.Port,
		User:   config.Database.User,
		Pass:   config.Database.Pass,
		DbName: config.Database.DbName,
	}
	err = pgsql.Init()
	if err != nil {
		logger.Panic("pgsql.Init", err)
	}

	err = api.Init()
	if err != nil {
		logger.Panic("api.Init", err)
	}
}
