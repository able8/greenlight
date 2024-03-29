package main

import (
	"context"
	"database/sql"
	"expvar"
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"

	// Import the pq driver so that it can register itself with the database/sql package.
	"github.com/able8/greenlight/internal/data"
	"github.com/able8/greenlight/internal/jsonlog"
	"github.com/able8/greenlight/internal/mailer"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

// Declare a string containing the application version number.
// const version = "1.0.0"

// Create a buildTime variable to hold the executable binary build time.
// Note that this must bu a string type, as the -X linker flag will only work
// with string variables.
var (
	buildTime string
	version   string
)

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
	// Add a new limiter struct containing fields for the requests-per-second and
	// burst values, and a boolean field which we can use to enable/disable rate limiting altogether.
	limiter struct {
		rps     float64
		burst   int
		enabled bool
	}
	// Update the config struct to hold the SMTP server settings.
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	// Add a cors struct and trustedOrigins field which the type []string.
	cors struct {
		trustedOrigins []string
	}
}

// Declare an application struct to hold the dependencies for out HTTP handlers, helpers, and middleware.
type application struct {
	config config
	// Change the logger field to have the type *jsonlog.Logger, instead of *log.Logger
	// logger *log.Logger
	logger *jsonlog.Logger
	models data.Models // Add a models struct to hold our new Models struct.
	mailer mailer.Mailer
	// Include a sync.WaitGroup in the application struct. The zero value for a sync.WaitGroup
	// type is a valid, useable, sync.WaitGroup with a counter value of 0,
	// so we don't need to do anything else to initialize it before we can use it.
	wg sync.WaitGroup
}

func main() {
	// Declare a instance of the config struct.
	var cfg config

	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")

	// Read the DSN value from the db-dsn command-line flags into the config struct.
	// We default to using our development DSN if no flag is provided.
	// flag.StringVar(&cfg.db.dsn, "db-dsn", "postgres://greenlight:password@localhost:/greenlight?sslmode=disable", "PostgresSQL DSN")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GREENLIGHT_DB_DSN"), "POSTGRESSQL DSN")

	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgresSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgresSQL max idle connections")
	flag.StringVar(&cfg.db.maxIdleTime, "db-max-idle-time", "15m", "PostgresSQL max connection idle time")

	// Create command line flags to read the settings values into the config struct.
	// Notice that we use true as the default for the 'enable' settings.
	flag.Float64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enable", true, "Enable rate limiter")

	// Read the SMTP server configuration settings into the config struct, using the
	// Mailtrap settings as the default values.
	flag.StringVar(&cfg.smtp.host, "smtp-host", "smtp.mailtrap.io", "SMTP host")
	flag.IntVar(&cfg.smtp.port, "smtp-port", 25, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", "f54ec4eee4356f", "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", "e031ff8accac48", "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", "Greenlight <no-replay@greenlight.com>", "SMTP sender")

	// Use the flag.Func() function to process the -cors-trusted-origins command line flag.
	// In this we use the strings.Fields() function to split the flag value into a
	// slice based on whitespace characters and assign it to our config struct.
	// Importantly, if the -cors-trusted-origins flag is not present, contains the empty
	// string, or contains only whitespace, then strings.Fields() will return an empty []string slice.
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	// Create a new version boolean flag with the default value of false.
	displayVersion := flag.Bool("version", false, "Display version and exit")

	flag.Parse()

	// If the version flag value is true, then print out the version number and exit.
	if *displayVersion {
		fmt.Printf("Version:\t%s\n", version)
		fmt.Printf("Build time:\t%s\n", buildTime)
		os.Exit(0)
	}

	// Initialize a new logger which writes messages to the standard out stream,
	// prefixed with the current date and time.
	// logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	// Initialize a new jsonlog.Logger which writes any messages at or above
	// the INFO severity level to the standard out stream.
	logger := jsonlog.New(os.Stdout, jsonlog.LevelInfo)

	// Call the openDB() helper function to create the connection pool,
	// passing in the config struct. If this returns an error,
	// we log it and exit the application immediately.
	db, err := openDB(cfg)
	if err != nil {
		// 	use the printfatal() method to write a log entry containing the error at the
		// 	fatal level and exit.  we have no additional properties to include in the log
		// 	entry, so we pass nil as the second parameter.
		logger.PrintFatal(err, nil)

	}
	// Defer a call to db.Close() so that the connection pool is closed be fore the main() exits.
	defer db.Close()

	// Also log a message to say that the connection pool has been successfully established.
	logger.PrintInfo("database connection pool established", nil)

	// migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	// if err != nil {
	// 	logger.PrintFatal(err, nil)
	// }

	// // migrator, err := migrate.NewWithDatabaseInstance("file:migrations/", "postgres", migrationDriver)
	// // if err != nil {
	// // 	logger.PrintFatal(err, nil)
	// // }

	// // err = migrator.Up()
	// // if err != nil && err != migrate.ErrNoChange {
	// // 	logger.PrintFatal(err, nil)
	// // }

	// // logger.PrintInfo("database migrations applied", nil)

	// Publish a new "version" variable in the expvar handler containing version number.
	expvar.NewString("version").Set(version)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() interface{} {
		return runtime.NumGoroutine()
	}))

	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() interface{} {
		return db.Stats()
	}))

	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() interface{} {
		return time.Now().Unix()
	}))

	// Read the value of the port and env command-line flags into the config struct.
	// Declare an instance of the application struct,
	// containing the config struct and the logger.
	app := &application{
		config: cfg,
		logger: logger,
		// User the data.NewModels() function to initialize a Models struct, passing
		// in the connection pool as a parameter.
		models: data.NewModels(db),
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		logger.PrintFatal(err, nil)
	}

	// Declare a new servemux and add a route.
	// mux := http.NewServeMux()
	// mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)
	// Use the httprouter instance returned by app.routes() as the server handler.

	// Declare a HTTP server with some sensible timeout settings.
	// srv := &http.Server{
	// 	Addr:    fmt.Sprintf(":%d", cfg.port),
	// 	Handler: app.routes(),
	// 	// Create a new Go log.Logger instance, passing in
	// 	// our custom Logger as the first parameter. The "" and 0 indicate
	// 	// that the log.Logger instance should not use a prefix or any flags.
	// 	ErrorLog:     log.New(logger, "", 0),
	// 	IdleTimeout:  time.Minute,
	// 	ReadTimeout:  10 * time.Second,
	// 	WriteTimeout: 30 * time.Second,
	// }

	// Start the HTTP server.
	// logger.Printf("Starting %s server on %s", cfg.env, srv.Addr)

	// logger.PrintInfo("Starting %s server on %s", map[string]string{
	// 	"addr": srv.Addr,
	// 	"env":  cfg.env,
	// })

	// Because the err variable is now already declared in the code above,
	// we need to use the = operator here, instead of the := operator.
	// err = srv.ListenAndServe()
	// logger.Fatal(err)
	// logger.PrintFatal(err, nil)
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
