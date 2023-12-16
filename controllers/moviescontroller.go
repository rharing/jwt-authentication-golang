package controllers

import (
	"github.com/gin-gonic/gin"
	"jwt-authentication-golang/movies"
	"net/http"
)

func GetCities(context *gin.Context) {
	cities, err := movies.LocateCities("")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {
		context.JSON(http.StatusOK, gin.H{"cities": cities})
	}
}
