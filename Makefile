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
	go run ./cmd/api

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
