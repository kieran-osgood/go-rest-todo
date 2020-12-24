include .env
export $(shell sed 's/=.*//' .env)

up: 
	docker-compose --env-file .env up --no-recreate 
down: 
	docker-compose down --rmi all
run:
	go run ./src/gqlgen/server.go
generate:
	go run github.com/99designs/gqlgen generate
db-generate: # usage: make db-generate migration="INSERT_MIGRATION_NAME"
	migrate create -ext sql -dir database/migrations -seq $(migration)
db-migrate:
	migrate -path database/migrations -database postgres://$$HOST:$$PORT/database up 1
