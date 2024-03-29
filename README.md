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

// Return the 5 records on page 1 (records 1-5 in the dataset)
curl "http://localhost:4000/v1/movies?page=1&page_size=5"

// Return the next 5 records on page 2 (records 6-10 in the dataset)
curl "http://localhost:4000/v1/movies?page=2&page_size=5"

Behind the scenes, the simplest way to support this style of pagination is by addingLIMIT and OFFSET clauses to our SQL query.

### 9.8. Returning Pagination Metadata

```sh
curl "http://localhost:4000/v1/movies?page=2&page_size=1"

{
        "metadata": {
                "current_page": 2,
                "page_size": 1,
                "first_page": 1,
                "last_page": 4,
                "total_records": 4
        },
        "movie": [
                {
                        "id": 2,
                        "title": "Black Panther",
                        "year": 2018,
                        "runtime": "134 mins",
                        "genres": [
                                "sci-fi",
                                "action",
                                "adventure"
                        ],
                        "version": 2
                }
        ]
}
```

## 10. Structured Logging and Error Handling

### 10.1. Structured JSON Log Entries

```sh
mkdir internal/jsonlog
touch internal/jsonlog/jsonlog.go

greenlight git:(main) ✗ go run ./cmd/api

{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"database connection pool established"}
{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"database migrations applied"}
{"level":"INFO","time":"2021-06-23T15:06:37Z","message":"Starting %s server on %s","properties":{"addr":":4000","env":"development"}}
{"level":"FATAL","time":"2021-06-23T15:06:37Z","message":"listen tcp :4000: bind: address already in use"}
```

### 10.2. Panic Recovery

## 11. Rate Limiting

### 11.1. Global Rate Limiting

```sh
go get golang.org/x/time/rate@latest
touch cmd/api/middleware.go

for i in {1..6}; do 
curl "http://localhost:4000/v1/healthcheck"
done

{
    "error": "rate limit exceeded"
}
```

### 11.2. IP-based Rate Limiting

### 11.3. Configuring the Rate Limiters

```sh
go run ./cmd/api --limiter-burst=2

go run ./cmd/api --limiter-enable=false
```

## 12. Graceful Shutdown

### 12.1. Sending Shutdown Signals

```
pgrep -l api

pkill -SIGKILL api
pkill -SIGTERM api
```
you can also try sending a SIGQUIT signal - either by pressing Ctrl+\ on your keyboard or
running `pkill -SIGQUIT api`.

This will cause the application to exit with a stack dump, similar to this.


### 12.2. Intercepting Shutdown Signals

```
pkill -SIGKILL api
pkill -SIGTERM api
pkill -SIGINT api
```

### 12.3. Executing the Shutdown

Specifically, after receiving one of these signals we will call the Shutdown() method onour HTTP server. The official documentation describes this as follows:

> Shutdown gracefully shuts down the server without interrupting any active connections. Shutdown works by first closing all open listeners, then closing allidle connections, and then waiting indefinitely for connections to return to idle and then shut down.

```sh
➜  greenlight git:(main) ✗ curl -i localhost:4000/v1/healthcheck & pkill -SIGTERM api
[1] 88864
➜  greenlight git:(main) ✗ HTTP/1.1 200 OK
Content-Type: application/json
Date: Thu, 24 Jun 2021 13:32:19 GMT
Content-Length: 102
Connection: close

{
        "status": "available",
        "system_info": {
                "environment": "development",
                "version": "1.0.0"
        }
}

[1]  + 88864 done       curl -i localhost:4000/v1/healthcheck


{"level":"INFO","time":"2021-06-24T13:33:39Z","message":"Starting server","properties":{"addr":":4000","env":"development"}}
{"level":"INFO","time":"2021-06-24T13:33:41Z","message":"caught signal, shutting down server","properties":{"signal":"terminated"}}
{"level":"INFO","time":"2021-06-24T13:33:45Z","message":"stopped server","properties":{"addr":":4000"}}
```

## 13. User Model Setup and Registration

In the upcoming sections of this book, we’re going to shift our focus towards users:registering them, activating them, authenticating them, and restricting access to our APIendpoints depending on the permissions that they have.

### 13.1. Setting up the Users Database Table

```sh
migrate create -seq -ext=.sql -dir=./migrations create_users_table

CREATE TABLE IF NOT EXISTS users (
        id bigserial PRIMARY KEY,
        created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
        name text NOT NULL,
        email citext UNIQUE NOT NULL,
        password_hash bytea NOT NULL,
        activated bool NOT NULL,
        version integer NOT NULL DEFAULT 1
);

DROP TABLE IF EXISTS users;
```

The email column has the type citext (case-insensitive text). This type stores text dataexactly as it is inputted — without changing the case in any way — but comparisonsagainst the data are always case-insensitive… including lookups on associated indexes.

The password_hash column has the type bytea (binary string). In this column we’ll storea one-way hash of the user’s password generated using bcrypt — not the plaintextpassword itself.

```
Postgresql: ERROR: type “citext” does not exist

psql postgres

postgres=# \c greenlight
You are now connected to database "greenlight" as user "able".
greenlight=# create extension citext;
CREATE EXTENSION
greenlight=> \d users
```

The extension needs to be created in each database. If you want to automatically have an extension created, you can create it in the template1 database which (by default, at least) is the database used as a model for "create database", so with appropriate permissions, in psql:

```
\c template1
create extension citext;
```

Then new databases will include citext by default.


### 13.2. Setting up the Users Model

```
go get golang.org/x/crypto/bcrypt@latest
```
### 13.3. Registering a User

```
touch cmd/api/users.go

BODY='{"name": "Alice Smith", "email": "alice@example.com","password": "password"}'
curl -i -d "$BODY" localhost:4000/v1/users

BODY='{"name": "", "email": "aliceexample.com","password": "passw"}'
curl -i -d "$BODY" localhost:4000/v1/users

BODY='{"name": "Alice Smith", "email": "alice@example.com","password": "password"}'
curl -i -d "$BODY" localhost:4000/v1/users
```

## 14. Sending Emails

### 14.1. SMTP Server Setup

#### What is Mailtrap? https://mailtrap.io/

Mailtrap is a test mail server solution that allows testing email notifications without sending them to the real users of your application. Not only does Mailtrap work as a powerful email test tool, it also lets you view your dummy emails online, forward them to your regular mailbox, share with the team and more! Mailtrap is a mail server test tool built by Railsware Products, Inc., a premium software development consulting company.

### 14.2. Creating Email Templates

```
mkdir -p internal/mailer/templates
touch internal/mailer/templates/user_welcome.tmpl.html
```

### 14.3. Sending a Welcome Email

```
go get github.com/go-mail/mail/v2@v2.3.0

code internal/mailer/mailer.go
```

we’re also going to use the new Go 1.16 embedded filesfunctionality, so that the email template files will be built into our binary whenwe create it later. This is really nice because it means we won’t have to deploythese template files separately to our production server.

If you want toinclude these files you should use the * wildcard character in the path, like//go:embed "templates/*"You can specify multiple directories and files in one directive. For example://go:embed "images" "styles/css" "favicon.ico".

```
BODY='{"name": "Bob 3", "email": "bob3@example.com","password": "password"}'
curl -i -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users
```
### 14.4. Sending Background Emails

Recovering panics

It’s important to bear in mind that any panic which happens in this backgroundgoroutine will not be automatically recovered by our recoverPanic() middleware orGo’s http.Server, and will cause our whole application to terminate.

So we need to make sure that any panic in this background goroutine is manually recovered, using a similar pattern to the one in our recoverPanic() middleware.

If you need to execute a lot of background tasks in your application, it can get tediousto keep repeating the same panic recovery code — and there’s a risk that you mightforget to include it altogether.To help take care of this, it’s possible to create a simple helper function which wrapsthe panic recovery logic. 


```go
// The background() helper accepts an arbitrary function as a parameter.
func (app *application) background(fn func()) {
	// Launch a background goroutine.
	go func() {
		// Recover any panic
		defer func() {
			if err := recover(); err != nil {
				app.logger.PrintError(fmt.Errorf("%s", err), nil)
			}
		}()

		// Execute the function that we passed as the parameter.
		fn()
	}()
}
```

### 14.5. Graceful Shutdown of Background Tasks

When we initiate a graceful shutdown of our application, it won’t wait for any background goroutines that we’ve launched to complete.

Fortunately, we can prevent this by using Go’s sync.WaitGroup functionality tocoordinate the graceful shutdown and our background goroutines.


To try this out, go ahead and restart the API and then send a request to thePOST /v1/users endpoint immediately followed by a SIGTERM signal. For example:

```sh
BODY='{"name": "Bob 10", "email": "bob10@example.com","password": "password"}'
curl -i -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users & pkill -SIGTERM api &
```

```json
{"level":"INFO","time":"2021-06-25T16:22:02Z","message":"Starting server","properties":{"addr":":4000","env":"development"}}
{"level":"INFO","time":"2021-06-25T16:22:14Z","message":"caught signal, shutting down server","properties":{"signal":"terminated"}}
{"level":"INFO","time":"2021-06-25T16:22:14Z","message":"completing background tasks","properties":{"addr":":4000"}}
{"level":"INFO","time":"2021-06-25T16:22:19Z","message":"Send email successfully","properties":{"email":"bob10@example.com"}}
{"level":"INFO","time":"2021-06-25T16:22:19Z","message":"stopped server","properties":{"addr":":4000"}}
```

This nicely illustrates how the graceful shutdown process waited for the welcomeemail to be sent (which took about two seconds in my case) before finally terminatingthe application.

## 15. User Activation

### 15.1. Setting up the Tokens Database Table

```
migrate create -seq -ext=.sql -dir=./migrations create_tokens_table

CREATE TABLE IF NOT EXISTS tokens (
        hash bytea PRIMARY KEY,
        user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
        expiry timestamp(0) with time zone NOT NULL,
        scope text NOT NULL
);

DROP TABLE IF EXISTS tokens;
```

### 15.2. Creating Secure Activation Tokens

```
code internal/data/tokens.go
```

### 15.3. Sending Activation Tokens

```sh
BODY='{"name": "Bob 11", "email": "bob11@example.com","password": "password"}'
curl -i -w '\nTime: %{time_total}\n' -d "$BODY" localhost:4000/v1/users 
```

### 15.4. Activating a User

```
curl -X PUT -d '{"token": "invalid"}' localhost:4000/v1/users/activated


curl -X PUT -d '{"token": "MO7H725URI5OIDXRJERAB6HAKA"}' localhost:4000/v1/users/activated

psql $GREENLIGHT_DB_DSN
select email, activated, version FROM users;

```
## 16. Authentication

Remember: 

- Authentication is about confirming who a user is
- Authorization is about checking whether that user is permitted to dosomething

### 16.1. Authentication Options

### 16.2. Generating Authentication Tokens

The client sends a JSON request to a new POST/v1/tokens/authenticationendpoint containing their credentials (email and password).
We look up the user record based on the email, and check if the passwordprovided is the correct one for the user. 

If the password is correct, we use our app.models.Tokens.New() method togenerate a token with an expiry time of 24 hours and the scope"authentication".We send this authentication token back to the client in a JSON responsebody.

```
code cmd/api/tokens.go

BODY='{"email": "bob11@example.com","password": "password"}'
curl -i -d "$BODY" localhost:4000/v1/tokens/authentication 

{
        "authorization_token": {
                "token": "H4SAGFWIOXHTDFF7ZX2QLFTHDU",
                "expiry": "2021-06-27T11:24:54.408236+08:00"
        }
}
```

### 16.3. Authenticating Requests

Now that our clients have a way to exchange their credentials for anauthentication token, let’s look at how we can use that token to authenticatethem, so we know exactly which user a request is coming from.

Essentially, once a client has an authentication token we will expect them toinclude it with all subsequent requests in an Authorization header, like so:

```
Authorization: Bearer IEYZSSSSSSSSSSSSSSSSSS

curl localhost:4000/v1/healthcheck
curl -d '{"email": "bob11@example.com","password": "password"}' localhost:4000/v1/tokens/authentication 

curl -H "Authorization: Bearer MU2WSVFIUVKUGAUF77F3JOFF3A" localhost:4000/v1/healthcheck
curl -i -H "Authorization: Bearer xxxx" localhost:4000/v1/healthcheck
```

## 17. Permission-based Authorization

### 17.1. Requiring User Activation

```
curl -i localhost:4000/v1/movies/1

BODY='{"name": "Bob 15", "email": "bob15@example.com","password": "password"}'
curl -i -d "$BODY" localhost:4000/v1/users

curl -d '{"email": "bob15@example.com","password": "password"}' localhost:4000/v1/tokens/authentication

curl -i -H "Authorization: Bearer S4A6QP6JUWDLKOIP7TEUUUKQGA" localhost:4000/v1/movies/1

select email from users where activated=true;
curl -d '{"email": "bob13@example.com","password": "password"}' localhost:4000/v1/tokens/authentication

curl -i -H "Authorization: Bearer VCYXGNS7LYXI3NXXEV6SPLGO5E" localhost:4000/v1/movies/1
```

### 17.2. Setting up the Permissions Database Table

```
migrate create -seq -ext=.sql -dir=./migrations add_permissions

CREATE TABLE IF NOT EXISTS permissions (
        id bigserial PRIMARY KEY,
        code text NOT NULL
);

CREATE TABLE IF NOT EXISTS users_permissions (
        user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
        permission_id bigint NOT NULL REFERENCES permissions ON DELETE CASCADE,
        PRIMARY KEY (user_id, permission_id)
);

INSERT INTO permissions (code)
VALUES
        ('movies:read'),
        ('movies:write');


DROP TABLE IF EXISTS users_permissions;
DROP TABLE IF EXISTS permissions;
```

### 17.3. Setting up the Permissions Model

### 17.4. Checking Permissions

```
psql $GREENLIGHT_DB_DSN

# Set the activated field to true
UPDATE users SET activated = true WHERE email = 'alice@example.com';
UPDATE users SET activated = true WHERE email = 'bob15@example.com';

# Give all users the "movies:read" permission
INSERT INTO users_permissions
SELECT id, (SELECT id FROM permissions WHERE code ='movies:read') FROM users;

# Give bob15@example.com the "movies:write" permission
INSERT INTO users_permissions
VALUES(
        (SELECT id FROM users WHERE email = 'bob15@example.com'),
        (SELECT id FROM permissions WHERE code = 'movies:write')

);

# List all activated users and their permissions.
SELECT email, array_agg(permissions.code) as permissions
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.activated = true
GROUP BY email;

alice@example.com | {movies:read}
bob15@example.com | {movies:read,movies:write}

curl -d '{"email": "alice@example.com","password": "password"}' localhost:4000/v1/tokens/authentication

curl -i -H "Authorization: Bearer PZLA6XI2N2VRZNFQIEBPEH353M" localhost:4000/v1/movies/1

curl -X DELETE -H "Authorization: Bearer PZLA6XI2N2VRZNFQIEBPEH353M" localhost:4000/v1/movies/1
{
    "error": "your user account does not have the necessary permissions to access this resource"
}

curl -d '{"email": "bob15@example.com","password": "password"}' localhost:4000/v1/tokens/authentication
curl -i -H "Authorization: Bearer HTMKB472YTWCGS3BMMP4J5LSQA" localhost:4000/v1/movies/1

curl -X DELETE -H "Authorization: Bearer HTMKB472YTWCGS3BMMP4J5LSQA" localhost:4000/v1/movies/1
{
        "message": "movie successfully deleted"
}
```

### 17.5. Granting Permissions

```
BODY='{"name": "Bob 17", "email": "bob17@example.com","password": "password"}'
curl -i -d "$BODY" localhost:4000/v1/users

SELECT email, code FROM users
INNER JOIN users_permissions ON users.id = users_permissions.user_id
INNER JOIN permissions ON users_permissions.permission_id = permissions.id
WHERE users.email = 'bob17@example.com';
```

## 18. Cross Origin Requests (CORS)

We’re going to switch to a completely new topic and update our application so that it supports cross-origin requests (CORS) from JavaScript.

### 18.1. An Overview of CORS

A webpage on one origin can send data to a different origin. But a webpage on one origin is not allowed to receive data from a different origin.

The same-origin policy prevents a (potentially malicious) website on another origin from reading (possibly confidential) information from your website.

For example, if you have an API at api.example.com and a trusted JavaScript front-end application running on www.example.com, then you’ll probably want to allow cross-origin requests from the trusted www.example.com domain to your API.

Or perhaps you have a completely open public API, and you want to allow cross-originrequests from anywhere so it’s easy for other developers to integrate with their ownwebsites.


### 18.2. Demonstrating the Same-Origin Policy

```
mkdir -p cmd/examples/cors/simple
code cmd/examples/cors/simple/main.go

go run ./cmd/examples/cors/simple

localhost:9000
```

Cross-Origin Request Blocked: The Same Origin Policy disallows reading theremote resource at http://localhost:4000/v1/healthcheck.

The request was sent to our API, which processed the request and returned a successful 200 OK response to the webbrowser containing all our standard response headers. 
To re-iterate: the request itselfwas not prevented by the same-origin policy — it’s just that the browser won’t let JavaScript see the response.

Finally, it’s important to emphasize that the same-origin policy is a web browser thingonly.

Outside of a web browser, anyone can make a request to our API from anywhere,using curl, wget or any other means and read the response. That’s completely unaffected and unchanged by the same-origin policy.

### 18.3. Simple CORS Requests

Let’s now make some changes to our API which relax the same-origin policy, so that JavaScript can read the responses from our API endpoints. To start with, the simplest way to achieve this is by setting the following header on allour API responses: `Access-Control-Allow-Origin: *`

The Access-Control-Allow-Origin response header is used to indicate to a browser that it’s OK to share a response with a different origin. In this case, the header value isthe wildcard * character, which means that it’s OK to share the response with anyother origin.

```
go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001"

go run ./cmd/examples/cors/simple --addr=":9001"
go run ./cmd/examples/cors/simple --addr=":9002"
```

In order to process this command-line flag, we can combine the new Go 1.16 flags.Func() and strings.Fields() functions to split the origin values into a[]string slice ready for use.


### 18.4. Preflight CORS Requests

When a cross-origin request doesn’t meet these conditions, then the web browser will trigger an initial ‘preflight’ request before the real request. The purpose of this preflight request is to determine whether the real cross-origin request will be permitted or not.

```
mkdir -p cmd/examples/cors/preflight
code cmd/examples/cors/preflight/main.go

go run ./cmd/examples/cors/preflight
```

已拦截跨源请求：同源策略禁止读取位于 http://localhost:4000/v1/tokens/authentication 的远程资源。（原因：CORS 预检响应的 'Access-Control-Allow-Headers'，不允许使用头 'content-type'）。

Cross-Origin Request Blocked: The Same Origin Policy disallows reading theremote resource at http://localhost:4000/v1/tokens/authentication. (Reason:header ‘content-type’ is not allowed according to header ‘Access-Control-Allow-Headers’ from CORS preflight response).


An Access-Control-Allow-Headers header listing the request headers that can beincluded in real cross-origin requests to the URL.


##### Caching preflight responses

If you want, you can also add an Access-Control-Max-Age header to your preflight responses. This indicates the number of seconds that the information provided by the Access-Control-Allow-Methods and Access-Control-Allow-Headers headers can becached by the browser.


##### Preflight wildcards

Wildcards in these headers are currently only supported by 74% of browsers. Anybrowsers which don’t support them will block the preflight request.

The Authorization header cannot be wildcarded. Instead, you will need to includethis explicitly in the header like `Access-Control-Allow-Headers: Authorization, *`. Wildcards are not supported for credentialed requests (those with cookies or HTTP basic authentication). For these, the character * will be treated as the literal string"*", rather than as a wildcard.


## 19. Metrics

### 19.1. Exposing Metrics with Expvar

The expvar package is a public variable package provided by Golang, which can assist in debugging global variables. 

localhost:4000/debug/vars

### 19.2. Creating Custom Metrics

```
BODY='{"email": "bob17@example.com","password":"password"}'
#curl -i -d "$BODY" http://localhost:4000/v1/tokens/authentication
hey -d "$BODY" -m "POST" http://localhost:4000/v1/tokens/authentication

go run ./cmd/api -cors-trusted-origins="http://localhost:9000 http://localhost:9001" --limiter-enable=false
```

### 19.3. Request-level Metrics

In this chapter we’re going to create some new middleware to record custom request-level metrics for our application.

### 19.4. Recording HTTP Status Codes

```
go get github.com/felixge/httpsnoop@v1.0.1
```

## 20. Building, Versioning and Quality Control

We’re going shift our focus from writing code to managing and maintaining our project, and take steps to help automate common tasks and prepare our API for deployment.

### 20.1. Creating and Using Makefiles

A phony target is one that is not really the name of a file; rather it is just aname for a rule to be executed.

```Makefile
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
```

make db/migrations/new name=test

### 20.2. Managing Environment Variables

Using a .envrc file

```
touch .envrc

export GREENLIGHT_DB_DSN="postgres://greenlight:password@localhost:/greenlight?sslmode=disable"

echo '.envrc' >> .gitignore
```

### 20.3. Quality Controlling Code

Use the `go mod tidy` command to prune any unused dependencies from the go.modand go.sum files, and add any missing dependencies.

### 20.4. Module Proxies and Vendoring

```
go env
go env | grep -i proxy
GONOPROXY=""
GOPROXY="https://proxy.golang.org,direct"
```

But if you don’t want to use the module mirror provided by Google, or you’re behind afirewall that blocks it, there are other alternatives like https://goproxy.io and theMicrosoft-provided https://athens.azurefd.net that you can try instead. Or you caneven host your own module mirror using the open-source Athens and goproxy projects.

But if you don’t want to use the module mirror provided by Google, or you’re behind afirewall that blocks it, there are other alternatives like https://goproxy.io and the Microsoft-provided https://athens.azurefd.net that you can try instead. Or you can even host your own module mirror using the open-source Athens and goproxy projects.

#### Vendoring

proxy.golang.org does not save all modules forever. There are a number ofreasons for this, but one reason is if proxy.golang.org is not able to detect asuitable license. In this case, only a temporarily cached copy of the module willbe made available, and may become unavailable if it is removed from theoriginal source and becomes outdated.

So, for these reasons, it can still be sensible to vendor your project dependenciesusing the go mod vendor command. Vendoring dependencies in this way basicallystores a complete copy of the source code for third-party packages in a vendor folderin your project.

```
make vendor

➜  greenlight git:(main) ✗ tree -L 3 ./vendor 
./vendor
├── github.com
│   ├── felixge
│   │   └── httpsnoop
│   ├── go-mail
│   │   └── mail
│   ├── golang-migrate
│   │   └── migrate
│   ├── hashicorp
│   │   ├── errwrap
│   │   └── go-multierror
│   ├── julienschmidt
│   │   └── httprouter
│   └── lib
│       └── pq
├── golang.org
│   └── x
│       ├── crypto
│       └── time
├── gopkg.in
│   └── alexcesaro
│       └── quotedprintable.v3
└── modules.txt

21 directories, 1 file
```

Because all the dependency source code is now stored in your project repository itself,it’s easy to check it into Git (or an alternative version control system) alongside the restof your code. This is reassuring because it gives you complete ownership of all the code used to build and run your applications, kept under version control.

#### Vendoring new dependencies

Fortunately, like most other reverse proxies, Caddy adds an X-Forwarded-For headerto each request. This header will contain the real IP address for the client.

```
go get github.com/tomasen/realip@latest

make vendor
make run/api 
```


### 20.5. Building Binaries

```Makefile
## build/api: build the cmd/api application
.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	go build -o bin/api ./cmd/api
```

```sh
make build/api

go tool dist list
make build/api
tree ./bin

echo 'bin/' >> .gitignore
cat .gitignore

go env GOCACHE
```

### 20.6. Managing and Automating Version Numbers

```sh
make build/api

./bin/api --version

git tag v1.0.0

git describe --always --dirty --tags --long
v1.0.0-0-ge5770e8-dirty
```

```Makefile
## build/api: build the cmd/api application
current_time = $(shell date +%Y-%m-%dT%H:%M:%S%z)
git_description = $(shell git describe --always --dirty --tags --long)
linker_flags = '-s -X main.buildTime=${current_time} -X main.version=${git_description}'

.PHONY: build/api
build/api:
	@echo "Building cmd/api..."
	go build -ldflags=${linker_flags} -o bin/api ./cmd/api
	GOOS=linux GOARCH=amd64 go build -ldflags=${linker_flags} -o bin/linux_amd64/api ./cmd/api

```

## 21. Deployment and Hosting

- Automate the configuration of the server — including creating user accounts, configuring the firewall and installing necessary software.
- Automate the process of updating your application and deploying changes to the server.
- How to run your application as a background service using systemd, as a non-root user.
- Use Caddy as a reverse proxy in front of your application to automatically manage TLS certificates and handle HTTPS connections.


### 21.1. Creating a Digital Ocean Droplet

### 21.2. Server Configuration and Installing Software

```sh
mkdir -p remote/setup
code remote/setup/01.sh

# Add the new user (and give them sudo privileges).
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"

# Force a password to be set for the new user the first time they log in
passwd --delete "${USERNAME}"
chage --lastday 0 "${USERNAME}"

# Copy the SSH keys from the root user to the new User.
rsync --archive --chown=${USERNAME}:${USERNAME} /root/.ssh /home/${USERNAME}

# Install the migrate CLI tool.
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
mv migrate.linux-amd64 /usr/local/bin/migrate

# Install PostgreSQL
apt --yes install postgresql

# Set up the greenlight DB and create a user account with the password entered earlier.
sudo -i -u postgresql psql -c "CREATE DATABASE greenlight"
sudo -i -u postgresql psql -d greenlight -c "CREATE EXTENSION IF NOT EXISTS citext"
sudo -i -u postgresql psql -d greenlight -c "CREATE ROLE greenlight WITH LOGIN PASSWORD '${DB_PASSWORD}'"

sudo systemctl status caddy
```

```
rsync -rP --delete ./remote/setup root@xxx:/root
```

- `-r` flag indicates that we want to copy the contents of ./remote/setup recursively,
- `-P` flag indicates that we want to display progress of the transfer
- `--delete` flag indicates that we want to delete any extraneous files from destination directory

#### Future changes to the server configuration

If you need to make further changes to your droplet configuration or settings, you can create an additional `remote/setup/02.sh` script and then execute it in the following way:

```sh
rsync -rP --delete ./remote/setup greenlight@xxx:~
ssh -t greenlight@xxx "sudo bash /home/greenlight/setup/02.sh"
```


### 21.3. Deployment and Executing Migrations

At a very high-level, our deployment process will consist of three actions:

- Copying the application binary and SQL migration files to the droplet.
- Executing the migrations against the PostgreSQL database on the droplet.
- Starting the application binary as a background service.


### 21.4. Running the API as a Background Service

In order to run our API application as a background service, the first thing we need to do is make a unit file, which informs systemd how and when to run theservice.

```sh
mkdir remote/production
code remete/production/api.service

make production/configure/api
make production/connect
ps -U greenlight

sudo ufw delete allow 4000/tcp
sudo ufw status

# Viewing logs
sudo journal -u api -n 100 -r
```

To ‘install’ the file, we need to copy it into the /etc/systemd/system/ folder on our droplet.

Then we need to run the `systemctl enable api` command on our droplet tomake systemd aware of the new unit file and automatically enable the servicewhen the droplet is rebooted.

Finally, we need to run `systemctl restart api` to start the service.

### 21.5. Using Caddy as a Reverse Proxy

The simplest way to configure Caddy is to create a Caddyfile — which containsa series of rules describing what we want Caddy to do.

You can open a SSH tunnel to the droplet and view them usinga web browser on your local machine.

For example, you could open an SSHtunnel between port 4000 on the droplet and port 9999 on your local machine by running the following command (make sure to replace both IP addresseswith your own droplet IP).

Caddy will automatically handle provisioning and renewing TLS certificates foryour domain via Let’s Encrypt, as well as redirecting all HTTP requests toHTTPS. It’s simple to set up, very robust, and saves you the overhead of needing to keep track of certificate renewals manually.

For the final time, deploy this Caddyfile update to your droplet…

```sh
code remote/production/Caddyfile

make production/configure/caddyfile

ssh -L :9999:xxxx:4000 greenlight@xxxx

curl -i https://greenlight.xx.com
```

## 22. Appendices

### 22.1. Managing Password Resets

### 22.2. Creating Additional Activation Tokens

### 22.3. Authentication with JSON Web Tokens

### 22.4. JSON Encoding Nuances

### 22.5. JSON Decoding Nuances

### 22.6. Request Context Timeouts

