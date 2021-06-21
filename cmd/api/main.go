package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	// Import the pq driver so that it can register itself with the database/sql package.
	"github.com/able8/greenlight/internal/data"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Declare a string containing the application version number.
const version = "1.0.0"

// Declare a config struct to hold all the configuration settings for our application.

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

// Declare an application struct to hold the dependencies for out HTTP handlers, helpers, and middleware.
type application struct {
	config config
	logger *log.Logger
	models data.Models // Add a models struct to hold our new Models struct.
}

func main() {
	// Declare a instance of the config struct.
	var cfg config

	// Read the value of the port and env command-line flags into the config struct.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Read the DSN value from the db-dsn command-line flags into the config struct.
	// We default to using our development DSN if no flag is provided.
	flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost:/greenlight?sslmode=disable", "PostgresSQL DSN")

	// flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "POSTGRESSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSQL max connection idle time")

	flag.Parse()

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Call the openDB() helper function to create the connection pool,
	// passing in the config struct. If this returns an error,
	// we log it and exit the application immediately.
	db, err := openDB(cfg)
	if err != nil {
		logger.Fatal(err)
	}
	// Defer a call to db.Close() so that the connection pool is closed be fore the main() exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully established.
	logger.Printf("database connection pool established")

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal(err, nil)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file:migrations/", "postgres", migrationDriver)
	if err != nil {
		logger.Fatal(err, nil)
	}

	err = migrator.Up()
	if err != nil && err != migrate.ErrNoChange {
		logger.Fatal(err, nil)
	}

	logger.Printf("database migrations applied")

	// Declare an instance of the application struct,
	// containing the config struct and the logger.
	app := &application{
		config: cfg,
		logger: logger,
		// User the data.NewModels() function to initialize a Models struct, passing
		// in the connection pool as a parameter.
		models: data.NewModels(db),
	}

	// Declare a new servemux and add a route.
	// mux := http.NewServeMux()
	// mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	// Use the httprouter instance returned by app.routes() as the server handler.

	// Declare a HTTP server with some sensible timeout settings.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start the HTTP server.
	logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)

	// Because the err variable is now already declared in the code above,
	// we need to use the = operator here, instead of the := operator.
	err = srv.ListenAndServe()
	logger.Fatal(err)
}

// The openDB() function returns a sql.DB connection pool.
func openDB(cfg config) (*sql.DB, error) {
	// User sql.Open() to create an empty connection pool, using the DSN from the config struct.
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool.
	db.SetMaxOpenConns(cfg.db.maxOpenConns)

	db.SetMaxIdleConns(cfg.db.maxIdleConns)

	// Convert the idle timeout duration string to a time.Duration type.
	duration, err := time.ParseDuration(cfg.db.maxIdleTime)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxIdleTime(duration)

	// Create a connect with a 5-second timeout deadline.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Establish a new connection to the database, passing in the context we
	// createed above as a parameter. If the connection couldn't be
	// established successfully within the 5 seconds deadline, then this will returns an error.
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// Return the sql.DB connection pool.
	return db, nil
}
