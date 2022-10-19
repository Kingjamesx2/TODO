//Filename: cmd/api/main.go

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// The application version nimber
const version = "1.0.0"

// The configuration settings
// The config struct -a set of complex port properties that specify the data type of the complex data type elements or the schema of the data
type config struct {
	port int
	env  string //development, staging, production, etc
}

// Dependency injection - the process of supplying a resource that a given piece of code requires.
type application struct {
	config config
	logger *log.Logger
}

func main() {
	var cfg config
	// read in the flags that are needed to populate our config
	// a flag is a predefined bit or bit sequence that holds a binary value.(not sure)
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development | staging | production )")
	// To parse -is where a string of commands – usually a program – is separated into more easily processed components, which are analyzed for correct syntax and then attached to tags that define each component.
	flag.Parse()
	//Create a logger - Logging is a means of tracking events that happen when some software runs.
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)
	//Create an instance of our applications struct
	app := &application{
		config: cfg,
		logger: logger,
	}

	//create our new servemux
	mux := http.NewServeMux()
	mux.HandleFunc("/v1/healthcheck", app.healthcheckHandler)

	// create our http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	// Start our server
	logger.Printf("starting %s server on %s", cfg.env, srv.Addr)
	err := srv.ListenAndServe()
	logger.Fatal(err)
}
