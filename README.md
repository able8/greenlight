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