package main

import (
	"backend/internal/models"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
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

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.cookie.name {
			claims := &claims{}
			refreshToken := cookie.Value

			// parse the token to get the claims
			_, err := jwt.ParseWithClaims(
				refreshToken,
				claims,
				func(t *jwt.Token) (interface{}, error) {
					return []byte(app.jwtSecret), nil
				},
			)
			if err != nil {
				app.errorJson(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			// get the user id from the token claims
			userId, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.errorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			user, err := app.userRepository.GetUserById(userId)
			if err != nil {
				app.errorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}

			jwtUser := jwtUser{
				id:        user.Id,
				firstName: user.FirstName,
				lastName:  user.LastName,
			}

			tokenPairs, err := app.auth.generateTokenPair(&jwtUser)
			if err != nil {
				app.errorJson(w, errors.New("error generating tokens"), http.StatusUnauthorized)
				return
			}

			http.SetCookie(w, app.auth.getRefreshCookie(tokenPairs.RefreshToken))
			app.writeJson(w, http.StatusOK, tokenPairs)
		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.getExpiredRefreshCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) movieCatalog(w http.ResponseWriter, r *http.Request) {
	movies, err := app.moviesRepository.GetAllMovies()
	if err != nil {
		app.errorJson(w, err)
		return
	}

	_ = app.writeJson(w, http.StatusOK, movies)
}

func (app *application) getMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieId, err := strconv.Atoi(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	movie, err := app.moviesRepository.GetMovieById(movieId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	_ = app.writeJson(w, http.StatusOK, movie)
}

// for admin
func (app *application) movieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieId, err := strconv.Atoi(id)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	movie, genres, err := app.moviesRepository.GetMovieByIdForEdit(movieId)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		movie,
		genres,
	}

	_ = app.writeJson(w, http.StatusOK, payload)
}

func (app *application) getAllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.genresRepository.GetAllGenres()
	if err != nil {
		app.errorJson(w, err)
		return
	}

	_ = app.writeJson(w, http.StatusOK, genres)
}

func (app *application) insertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie

	err := app.readJson(w, r, &movie)
	if err != nil {
		app.errorJson(w, err)
		return
	}

	response := JsonResponse{
		Error:   false,
		Message: "movie updated",
	}

	app.writeJson(w, http.StatusAccepted, response)
}
