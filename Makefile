up: 
	docker-compose up 
run:
	go run ./src/gqlgen/server.go
generate:
	go run github.com/99designs/gqlgen generate