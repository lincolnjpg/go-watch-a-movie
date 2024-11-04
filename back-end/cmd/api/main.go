package main

import (
	"backend/internal/repository"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

const port = 8080

type application struct {
	dsn              string
	domain           string
	moviesRepository repository.MoviesRepository
	genresRepository repository.GenreRepository
	userRepository   repository.UserRepository
	auth
	jwtSecret   string
	jwtIssuer   string
	jwtAudience string
	apiKey      string
}

func main() {
	var app application

	flag.StringVar(
		&app.dsn,
		"dsn",
		"host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5",
		"Postgres connection string",
	)
	flag.StringVar(&app.jwtSecret, "jwt-secret", "verysecret", "signing secret")
	flag.StringVar(&app.jwtIssuer, "jwt-issuer", "example.com", "signing issuer")
	flag.StringVar(&app.jwtAudience, "jwt-audience", "example.com", "signing audience")
	flag.StringVar(&app.domain, "domain", "example.com", "domain")
	flag.StringVar(&app.cookie.domain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.apiKey, "api-key", "f2f6c03fc893958775e650cf691663b5", "api key")
	flag.Parse()

	connection, err := app.connectToDb()
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()

	app.auth = auth{
		issuer:        app.jwtIssuer,
		audience:      app.jwtAudience,
		secret:        app.jwtSecret,
		tokenExpiry:   time.Minute * 15,
		refreshExpiry: time.Hour * 24,
		cookie: cookie{
			path:   "/",
			name:   "__Host-refresh_token",
			domain: app.cookie.domain,
		},
	}

	app.moviesRepository = &repository.PostgresMoviesRepository{Db: connection}
	app.userRepository = &repository.PostgresUserRepository{Db: connection}
	app.genresRepository = &repository.PostgresGenresRepository{Db: connection}

	log.Println("Starting application on port", port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", port), app.routes())
	if err != nil {
		log.Fatal(err)
	}
}
