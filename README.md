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

## 7. CRUD Operations

### 7.1. Setting up the Movie Model

### 7.2. Creating a New Movie

```sh
BODY='{"title":"Moana","year": 2016, "runtime": "107 mins", "genres": ["animation", "adventure"]}'
curl -i -d "$BODY" localhost:4000/v1/movies

BODY='{"title":"Black Panther","year": 2018, "runtime": "134 mins", "genres": ["action", "adventure"]}'
curl -i -d "$BODY" localhost:4000/v1/movies

BODY='{"title":"Deadpool","year": 2016, "runtime": "108 mins", "genres": ["action", "comedy"]}'
curl -i -d "$BODY" localhost:4000/v1/movies

BODY='{"title":"The Breakfash Club","year": 1986, "runtime": "96 mins", "genres": ["drame"]}'
curl -i -d "$BODY" localhost:4000/v1/movies


SELECT * FROM movies;
```

### 7.3. Fetching a Movie

### 7.4. Updating a Movie

```sh
curl localhost:4000/v1/movies/2

BODY='{"title":"Black Panther","year": 2018, "runtime": "134 mins", "genres": ["sci-fi","action", "adventure"]}'
curl -i -X PUT -d "$BODY" localhost:4000/v1/movies/2
```

### 7.5. Deleting a Movie

```sh
curl -X DELETE localhost:4000/v1/movies/5
```

## 8. Advanced CRUD Operations

### 8.1. Handling Partial Updates

The key thing to notice here is that pointers have the zero-value nil.
So — in theory — we could change the fields in our input struct to be pointers.Then to see if a client has provided a particular key/value pair in the JSON, we cansimply check whether the corresponding field in the input struct equals nil or not.

```sh
curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4
curl -X PATCH -d '{"year": 1985, "title":""}' localhost:4000/v1/movies/4
```

### 8.2. Optimistic Concurrency Control

```sh
curl -X PATCH -d '{"year": 1985}' localhost:4000/v1/movies/4 &
curl -X PATCH -d '{"year": 1986}' localhost:4000/v1/movies/4 &

{           
        "movie": {
		...
        }
}
{
        "error": "unable to update the record due to an edit conflict, please try again"
}
```

### 8.3. Managing SQL Query Timeouts

This feature can be useful when you have a SQL query that is taking longer to runthan expected. When this happens, it suggests a problem — either with thatparticular query or your database or application more generally — and you probablywant to cancel the query (in order to free up resources), log an error for furtherinvestigation, and return a 500 Internal Server Error response to the client.

We’ll update ourSQL query to return a pg_sleep(10) value, which will make PostgreSQL sleep for 10seconds before returning its result.

```sh
curl -w '\nTime: %{time_total}s \n' localhost:4000/v1/movies/4
```

## 9. Filtering, Sorting, and Pagination

### 9.1. Parsing Query String Parameters

```sh
curl "localhost:4000/v1/movies?title=godfather&genres=crime,drame&page=1&page_size=5&sort=year"

curl "localhost:4000/v1/movies"
```

### 9.2. Validating Query String Parameters

```sh
curl "localhost:4000/v1/movies?title=godfather&genres=crime,drame&page=1&page_size=5&sort=year"

curl "localhost:4000/v1/movies?page=-1&page_size=-1&sort=foo"
```

### 9.3. Listing Data

```sh
curl "localhost:4000/v1/movies"
```

### 9.4. Filtering Lists

// List all movies
/v1/movies

// List movies where the title is case-insensitive exact match
/v1/movies?title=balck+panther

// List movies where the genres includes 'adventure'
/v1/movies?genres=adventure

// List movies where the title is a case-insensitive exact match for 'moana' AND the genres includes both animation and adventure
/v1/movies?title=moana&genres=animation,adventure

curl "http://localhost:4000/v1/movies?genres=adventure"

curl "http://localhost:4000/v1/movies?title=moana&genres=animation,adventure"

> Note: The + symbol in the query strings above is URL-encoded space character. Alternatively you could use %20 instead.. either will work in the context of a query string.


### 9.5. Full-Text Search

// Return all movies where the title includes the case-insensitive word 'panther'
/v1/movies?title=panther

curl "http://localhost:4000/v1/movies?title=panther"
curl "http://localhost:4000/v1/movies?title=the+club"

- Adding indexes

```sh
migrate create -seq -ext .sql -dir ./migrations add_movies_indexes

CREATE INDEX IF NOT EXISTS movies_title_idx ON movies USING GIN (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS movies_genres_idx ON movies USING GIN (genres);

DROP INDEX IF EXISS movies_title_idx;
DROP INDEX IF EXISS movies_genres_idx;
```

### 9.6. Sorting Lists

// Sort the movies on the title field in ascending alphabetical order.
curl "http://localhost:4000/v1/movies?sort=title"

// Sort the movies on the year field in descending numerical order.
curl "http://localhost:4000/v1/movies?sort=-year"

### 9.7. Paginating Lists

### 9.8. Returning Pagination Metadata

## 10. Structured Logging and Error Handling

### 10.1. Structured JSON Log Entries

### 10.2. Panic Recovery

## 11. Rate Limiting

### 11.1. Global Rate Limiting

### 11.2. IP-based Rate Limiting

### 11.3. Configuring the Rate Limiters

## 12. Graceful Shutdown

### 12.1. Sending Shutdown Signals

### 12.2. Intercepting Shutdown Signals

### 12.3. Executing the Shutdown

## 13. User Model Setup and Registration

### 13.1. Setting up the Users Database Table

### 13.2. Setting up the Users Model

### 13.3. Registering a User

## 14. Sending Emails

### 14.1. SMTP Server Setup

### 14.2. Creating Email Templates

### 14.3. Sending a Welcome Email

### 14.4. Sending Background Emails

### 14.5. Graceful Shutdown of Background Tasks

## 15. User Activation

### 15.1. Setting up the Tokens Database Table

### 15.2. Creating Secure Activation Tokens

### 15.3. Sending Activation Tokens

### 15.4. Activating a User

## 16. Authentication

### 16.1. Authentication Options

### 16.2. Generating Authentication Tokens

### 16.3. Authenticating Requests

## 17. Permission-based Authorization

### 17.1. Requiring User Activation

### 17.2. Setting up the Permissions Database Table

### 17.3. Setting up the Permissions Model

### 17.4. Checking Permissions

### 17.5. Granting Permissions

## 18. Cross Origin Requests

### 18.1. An Overview of CORS

### 18.2. Demonstrating the Same-Origin Policy

### 18.3. Simple CORS Requests

### 18.4. Preflight CORS Requests

## 19. Metrics

### 19.1. Exposing Metrics with Expvar

### 19.2. Creating Custom Metrics

### 19.3. Request-level Metrics

### 19.4. Recording HTTP Status Codes

## 20. Building, Versioning and Quality Control

### 20.1. Creating and Using Makefiles

### 20.2. Managing Environment Variables

### 20.3. Quality Controlling Code

### 20.4. Module Proxies and Vendoring

### 20.5. Building Binaries

### 20.6. Managing and Automating Version Numbers

## 21. Deployment and Hosting

### 21.1. Creating a Digital Ocean Droplet

### 21.2. Server Configuration and Installing Software

### 21.3. Deployment and Executing Migrations

### 21.4. Running the API as a Background Service

### 21.5. Using Caddy as a Reverse Proxy

## 22. Appendices

### 22.1. Managing Password Resets

### 22.2. Creating Additional Activation Tokens

### 22.3. Authentication with JSON Web Tokens

### 22.4. JSON Encoding Nuances

### 22.5. JSON Decoding Nuances

### 22.6. Request Context Timeouts

