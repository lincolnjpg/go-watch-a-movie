package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCors)
	mux.Get("/home", app.Home)
	mux.Get("/movies", app.GetAllMovies)
	mux.Post("/authenticate", app.authenticate)

	return mux
}
