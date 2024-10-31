package main

import (
	"log"
	"net/http"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"status"`
		Message string `json:"message"`
		Version string `json:"version"`
	}{
		Status:  "active",
		Message: "Go Movies up and running",
		Version: "1.0.0",
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *application) GetAllMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.repository.GetAllMovies()
	if err != nil {
		app.errorJson(w, err)
		return
	}

	_ = app.writeJson(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read JSON payload

	// validate user against the db

	// check password

	// create a jwt user
	user := jwtUser{
		id:        1,
		firstName: "Alice",
		lastName:  "Smith",
	}

	// generate tokens
	tokens, err := app.auth.generateTokenPair(&user)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	log.Println(tokens.accesToken)
	refreshCookie := app.auth.getRefreshCookie(tokens.refreshToken)
	http.SetCookie(w, refreshCookie)
	w.Write([]byte(tokens.accesToken))
}
