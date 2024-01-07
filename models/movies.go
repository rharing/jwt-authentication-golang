package models

import (
	"time"
)

type Play struct {
	Movie      Movie     `json:"movie"`
	Start      time.Time `json:"start"`
	Tickethref string    `json:"ticket-href"`
}
type Cinema struct {
	Name  string `json:"Name"`
	Href  string `json:"href"`
	Plays []Play `json:"Plays"`
}
type Movie struct {
	Id        string `json:"id"`
	Title     string `json:"title"`
	Href      string `json:"href"`
	Rating    string `json:"rating"`
	Content   string `json:"content"`
	ImageHref string `json:"image-href"`
	Duration  int    `json:"duration"`
}
type City struct {
	Href    string    `json:"href"`
	Name    string    `json:"Name"`
	Cinemas []*Cinema `json:"Cinemas"`
}
type PlayDTO struct {
	Title      string    `json:"title"`
	MovieId    string    `json:"movie-id"`
	Moviehref  string    `json:"movie"`
	Tickethref string    `json:"ticket"`
	Start      time.Time `json:"start"`
}
type MoviesScraper interface {
	LocateCities(url string) ([]City, error)
	LocatePlaysForCity(city string) (City, error)
	LoadMovieContent(movieId string) (Movie, error)
}
type MyMovies struct {
	Wanted   []string
	Unwanted []string
	Seen     []string
}
type MoviesRepository interface {
	LoadMovieContent(id string) (Movie, error)
	SeenMovie(movieid string, userId string)
	WantedMovie(movieid string, userId string)
	UnwantedMovie(movieid string, userId string)
	ResetMovie(movieid string, userId string)
	MyMovies(userId string) MyMovies
}
