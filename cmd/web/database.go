package main

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

// initDB initializes the database connection
func initDB() *sql.DB {
	conn := connectToDB()
	if conn == nil {
		log.Panic("can't connect to the database...")
	}
	return conn
}

// connectToDB tries to make a connection in 10 seconds
func connectToDB() *sql.DB {
	counterToCancel := 0

	dsn := os.Getenv("DSN") // database address, get env.var. from Makefile

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("postgres not yet ready...")
		} else {
			log.Println("connected to the database!")
			return connection
		}

		if counterToCancel > 10 {
			return nil
		}
		log.Println("Backing off for 1 second...")
		time.Sleep(1 * time.Second)
		counterToCancel++
		continue
	}
}

// openDB opens a sql connection with the given database string
func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping() // ping the database
	if err != nil {
		return nil, err
	}

	return db, nil // successully connected
}
