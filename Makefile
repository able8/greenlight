# Include variables from the .envrc file
include .envrc

## help: print this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ":" | sed -e 's/^/  /'

.PHONY: confirm
confirm:
	@echo "Are you sure? [y/N]" && read ans && [ $${ans:-N} = y ]

## run/api: run the cmd/api application
.PHONY: run/api
# run/api: confirm
run/api: 
	@echo "Staring the application"
	go run ./cmd/api -db-dsn=${GREENLIGHT_DB_DSN}

## db/psql: connect to the database using psql
.PHONY: db/psql
db/psql:
	# suppress commands from being echoed by prefixing them with the @ character.
	@psql ${GREENLIGHT_DB_DSN}

## db/migrations/new name=$1: create a new database migration
.PHONY: db/migrations/new
db/migrations/new:
	echo "Creating migration files for ${name}..."
	migrate create -seq -ext=.sql -dir=./migrations ${name}

## vendor: tidy and vendor dependencies
.PHONY: vendor
vendor:
	@echo "Tidying and verifying module dependencies..."
	go mod tidy
	go mod verify
	@echo "Vendoring dependencies..."
	go mod vendor



## build/api: build the cmd/api application
current_time = $(shell date +%Y-%m-%dT%H:%M:%S%z)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'


.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	go build -ldflags=${linker_flags} -o bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o bin/linux_amd64/api ./cmd/api