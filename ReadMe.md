# Introduction
It's a Todo app - mindblown.

This is just an experiment to learn Go, with the aim to get a nice cohesive setup of packages.
# Get Started
In order to run the application you must have Docker and Docker-compose installed.

Prepare two separate terminals windows and run:
  1. `make up` - runs docker-compose to setup the postgresql database container and admin container for managing the database.
  2. `make air` - enables live reloading of the application
  (alternatively you can run `go get ./...` and then use `make run` after each save to view changes).

# Stack
`docker-compose` for setting up PostgresSQL/PGAdmin

[golang-migrate](https://github.com/golang-migrate/migrate) - automatic programmatic database migrations

[squirrel](https://github.com/Masterminds/squirrel) - composable Fluent SQL generation 

[zap](https://github.com/uber-go/zap) - structured performant error logging

[fuzzysearch](https://github.com/lithammer/fuzzysearch) - performant search suggestions.
