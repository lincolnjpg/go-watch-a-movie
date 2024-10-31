package main

import (
	"backend/internal/repository"
	"flag"
	"fmt"
	"log"
	"net/http"
)

const port = 8080

type application struct {
	Dsn        string
	Domain     string
	Repository repository.MoviesRepository
}

func main() {
	var app application

	flag.StringVar(
		&app.Dsn,
		"dsn",
		"host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5",
		"Postgres connection string",
	)
	flag.Parse()

	connection, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	app.Repository = &repository.PostgresMoviesRepository{Db: connection}
	app.Domain = "example.com"

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
