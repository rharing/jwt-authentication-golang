package controllers

import (
	"github.com/gin-gonic/gin"
	models "jwt-authentication-golang/models"
	"jwt-authentication-golang/movies"

	"net/http"
)

func GetCities(context *gin.Context) {
	cities, err := movies.LocateCities("")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		if err == nil {
			context.JSON(http.StatusOK, gin.H{"cities": cities})
		}
	}
}

func GetCity(context *gin.Context) {
	city := context.Param("city")
	CityWithPlays, err := movies.LocatePlaysForCity(city)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		if err == nil {
			playDtos := make([]models.PlayDTO, 0)
			for i := 0; i < len(CityWithPlays.Cinemas); i++ {
				cinema := CityWithPlays.Cinemas[i]
				for j := 0; j < len(cinema.Plays); j++ {
					play := cinema.Plays[j]
					playDtos = append(playDtos, models.PlayDTO{Title: play.Movie.Title, Href: play.Movie.Href, Start: play.Start})
				}
			}
			context.JSON(http.StatusOK, gin.H{"plays": playDtos})
		}
	}
}
