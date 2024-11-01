package main

import (
	"errors"
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
	movies, err := app.moviesRepository.GetAllMovies()
	if err != nil {
		app.errorJson(w, err)
		return
	}

	_ = app.writeJson(w, http.StatusOK, movies)
}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {
	// read JSON payload
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJson(w, r, &requestPayload)
	if err != nil {
		app.errorJson(w, err)

		return
	}

	// validate user against the db
	user, err := app.userRepository.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.errorJson(w, errors.New("invalid information"))

		return
	}

	// check password
	isPasswordValid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !isPasswordValid {
		app.errorJson(w, errors.New("invalid credentials"))
	}

	// create a jwt user
	jwtUser := jwtUser{
		id:        user.Id,
		firstName: user.FirstName,
		lastName:  user.LastName,
	}

	// generate tokens
	tokens, err := app.auth.generateTokenPair(&jwtUser)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	refreshCookie := app.auth.getRefreshCookie(tokens.RefreshToken)
	http.SetCookie(w, refreshCookie)
	app.writeJson(w, http.StatusAccepted, tokens)
}
