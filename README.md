# greenlight

---
《Let's Go Further: Advanced patterns for building APIs and web applications in Go 》

https://lets-go-further.alexedwards.net/

Contents: https://lets-go-further.alexedwards.net/sample/00.01-contents.html

---
《Let’s Go: Learn to build professional web applications with Go》

https://lets-go.alexedwards.net/

Contents: https://lets-go.alexedwards.net/sample/00.01-contents.html


## Introduction

In this book we’re going to work through the start-to-finish build of an application called Greenlight — a JSON API for retrieving and managing information about movies. You can think of the core functionality as being a 
bit like the Open Movie Database API.

## Chapter 2.1 Project Setup and Skeleton Structure

```bash
mkdir greenlight
go mod init github.com/able8/greenlight

mkdir -p bin cmd/api internal migrations remote
touch Makefile
touch cmd/api/main.go

go run ./cmd/api  
```

## Chapter 2.2 Basic HTTP Server

```bsh
go run ./cmd/api  

curl -i localhost:4000/v1/healthcheck

go run ./cmd/api  -port=3000 -env=production
```

http://localhost:4000/v1/healthcheck


## 2.3. API Endpoints and RESTful Routing

- Choosing a router

```bash
go get github.com/julienschmidt/httprouter@v1.3.0
```

```bash
curl -i localhost:4000/v1/healthcheck

curl -X POST localhost:4000/v1/movies

curl localhost:4000/v1/movies/123
```

## Chapter 3 Sending JSON Responses

## Chapter 3.1 Fixed-Format JSON

## Chapter 3.2 JSON Encoding

## Chapter 3.3 Encoding Structs

## Chapter 3.4 Formatting and Enveloping Responses

## Chapter 3.5 Advanced JSON Customization

## Chapter 3.6 Sending Error Messages

## Chapter 4 Parsing JSON Requests

## Chapter 4.1 JSON Decoding

```bash
# Create a BODY variable containing the JSON data that we want to send.
BODY='{"title":"Moana","year": 2016, "runtime": 107, "genres": ["animation", "adventure"]}'
# Use the -d flag to send the contents of the BODY variable as the HTTP request body.
curl -i -d "$BODY" localhost:4000/v1/movies
```

## Chapter 4.2 Managing Bad Requests

```bash
# Send a numeric 'title' value instead of string
curl -d '{"title": 123}' localhost:4000/v1/movies

# Send an empty request body
curl -X POST localhost:4000/v1/movies
```

## Chapter 4.3 Restricting Inputs

```bash
# Body contains multiple JSON values
curl -i -d '{"title": "Moana"}{"title": "Top Gun"}' localhost:4000/v1/movies
# Body contains garbage contents after the first JSON value
curl -i -d '{"title": "Moano"} :~()' localhost:4000/v1/movies
```

## Chapter 4.4.Custom JSON Decoding

```bash
curl -i -d '{"title": "Moana","runtime": "107 mins"}' localhost:4000/v1/movies

curl -i -d '{"title": "Moana","runtime": "107 minutes"}' localhost:4000/v1/movies
```

## Chapter 4.5.Validating JSON Input

```bash
mkdir internal/validator
touch internal/validator/validator.go

BODY='{"title":"","year":1000,"runtime":"-123 mins","genres":["sci-fi","sci-fi"]}'

curl -i -d "$BODY" localhost:4000/v1/movies
```

## Chapter 5 Database Setup and Configuration

## Chapter 5.1.Setting up PostgreSQL

```bash
brew install postgresql
psql --version
# https://wiki.postgresql.org/wiki/Homebrew
brew services start postgresql
psql postgres

SELECT current_user;
# Creating database, users, and extensions
CREATE DATABASE greenlight;
\c greenlight

CREATE ROLE greenlight WITH LOGIN PASSWORD 'passoword';
CREATE EXTENSION IF NOT EXISTS citext;
\q
```

- Connecting as the new user

```bash
psql --host=localhost --dbname=greenlight --username=greenlight

SELECT current_user;
\q

psql postgres -c 'SHOW config_file;'
```

## Chapter 5.2 Connecting to PostgreSQL

```sh
go get github.com/lib/pq@v1.10.0

export GREENLIGHT_DB_DSN="postgres://greenlight:password@localhost:/greenlight?sslmode=disable"
```

## Chapter 5.3 Configuring the Database Connection Pool


## Chapter 6 SQL Migrations

To do this, we could simply use the psql tool again and run the necessary CREATE TABLE statement against our database.

But instead, we’re going to explore howto use SQL migrations to create the table (and more generally, manage databaseschema changes throughout the project).

- Installing the migrate tool

```sh
brew install golang-migrate
migrate -veriosn
```


## Chapter 6.2.Working with SQL Migrations

```sh
migrate create -seq -ext=.sql -dir=./migrations create_movies_table

migrate create -seq -ext=.sql -dir=./migrations add_movies_check_constraints

# Executing the migrations
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN up

# Fixing errors in SQL Migrations
migrate -path=./migrations -database=$GREENLIGHT_DB_DSN force xx


psql  $GREENLIGHT_DB_DSN
\dt

select * from schema_migrations;

\d movies;

```
