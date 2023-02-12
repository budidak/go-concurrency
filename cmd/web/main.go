package main

import (
	"fmt"
	"go-concurrency/data"
	"log"
	"net/http"
	"os"
	"sync"
)

const webPort = "80"

func main() {
	// connect to the database
	db := initDB()

	// create sessions for server side rendering
	session := initSession()

	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// create channels

	// create wait group
	wg := sync.WaitGroup{}

	// set up application config
	app := Config{
		Session:  session,
		DB:       db,
		InfoLog:  infoLog,
		ErrorLog: errorLog,
		Wait:     &wg,
		Models:   data.New(db), // creates empty db models
	}

	// set up mail

	// start goroutine in the background for graceful shutdown
	go app.listenForShutdown()

	// listen for web connections
	app.serve()
}

// serve starts an http server and listen a port
func (app *Config) serve() {
	srv := &http.Server{
		Addr:    fmt.Sprintf("localhost:%s", webPort),
		Handler: app.routes(),
	}

	app.InfoLog.Println("starting web server...")
	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
