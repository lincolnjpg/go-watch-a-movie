package main

import (
	"backend/internal/models"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
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

	out, err := json.Marshal(payload)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func (app *application) AllMovies(w http.ResponseWriter, r *http.Request) {
	var movies []models.Movie

	releaseDate, _ := time.Parse("2006-01-02", "1986-03-07")

	highlander := models.Movie{
		Id:          1,
		Title:       "Highlander",
		ReleaseDate: releaseDate,
		MpaaRating:  "R",
		RunTime:     116,
		Description: "A very nice movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	movies = append(movies, highlander)

	releaseDate, _ = time.Parse("2006-01-02", "1981-06-12")

	raidersOfTheLostArk := models.Movie{
		Id:          2,
		Title:       "Raiders of the Lost Ark",
		ReleaseDate: releaseDate,
		MpaaRating:  "PG-13",
		RunTime:     115,
		Description: "Another very nice movie",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	movies = append(movies, raidersOfTheLostArk)

	out, err := json.Marshal(movies)
	if err != nil {
		fmt.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}
