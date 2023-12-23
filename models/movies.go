package models

import "time"

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
	Title string    `json:"title"`
	Href  string    `json:"href"`
	Start time.Time `json:"start"`
}
type MoviesScraper interface {
	LocateCities(url string) ([]City, error)
}
