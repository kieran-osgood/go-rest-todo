# Loads the backend .env file
include backend/.env
export $(shell sed 's/=.*//' backend/.env)

# BACKEND
up:
	cd backend && docker-compose --env-file .env up --no-recreate 
down:
	cd backend && docker-compose down --rmi all
run:
	cd backend && go run main.go
generate:
	cd backend && go run github.com/99designs/gqlgen generate
db-generate: # usage: make db-generate migration="INSERT_MIGRATION_NAME"
	cd backend && migrate create -ext sql -dir database/migrations -seq $(migration)
db-migrate:
	cd backend && migrate -path database/migrations -database postgres://$$HOST:$$PORT/database up 1
air:
	cd backend && air

# FRONTEND
dev:
	cd frontend && npm run dev