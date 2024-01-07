package controllers

import (
	"github.com/gin-gonic/gin"
	models "jwt-authentication-golang/models"
	"jwt-authentication-golang/movies"

	"net/http"
)

type RouteController struct {
}

var moviesProvider	models.MoviesRepository

func init() {
	moviesProvider = movies.NewInMemoryMoviesRepository()
}
const overviewLocation = "file://./movies/resources/overview_haarlem.html"
const oppenheimerLocation = "file://./movies/resources/oppenheimer.html"

func GetCities(context *gin.Context) {
	var url = "http://www.filmladder.nl"
	if context.Query("use4Testing") == "1" {
		url = overviewLocation
	}
	cities, err := movies.LocateCities(url)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		if err == nil {
			context.JSON(http.StatusOK, gin.H{"cities": cities})
		}
	}
}

func GetCity(context *gin.Context) {
	CityWithPlays, err := locatePlays(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		if err == nil {
			playDtos := make([]models.PlayDTO, 0)
			for i := 0; i < len(CityWithPlays.Cinemas); i++ {
				cinema := CityWithPlays.Cinemas[i]
				for j := 0; j < len(cinema.Plays); j++ {
					play := cinema.Plays[j]
					playDtos = append(playDtos, models.PlayDTO{
						Title:      play.Movie.Title,
						MovieId:    play.Movie.Id,
						Tickethref: play.Tickethref,
						Moviehref:  play.Movie.Href,
						Start:      play.Start})
				}
			}
			context.JSON(http.StatusOK, gin.H{"plays": playDtos})
		}
	}
}
func GetMovie(context *gin.Context) {
	movie, err := locateMovie(context)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		context.JSON(http.StatusOK, gin.H{"movie": movie})
	}
}

func locateMovie(context *gin.Context) (models.Movie, error) {
	var id = context.Param("movieId")
	url := locateExternalUrlFromId(id)
	if context.Query("use4Testing") == "1" {
		url = oppenheimerLocation
	}
	return models.MoviesRepository.LoadMovieContent(url)
}

func locateExternalUrlFromId(id string) string {

}

}

func locatePlays(context *gin.Context) (models.City, error) {
	var url = context.Param("city")

	if context.Query("use4Testing") == "1" {
		url = overviewLocation
		return movies.LocatePlays(url)
	}
	{
		return movies.LocatePlaysForCity(url)
	}
}
