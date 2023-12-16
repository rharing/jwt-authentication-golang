package controllers

import (
	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"
	"github.com/json-iterator/go/extra"
	"jwt-authentication-golang/movies"
	"net/http"
)

func GetCities(context *gin.Context) {
	cities, err := movies.LocateCities("")
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": err})
	} else {

		extra.SupportPrivateFields()
		jsonCities, err := jsoniter.ConfigCompatibleWithStandardLibrary.MarshalToString(cities)
		if err == nil {
			context.JSON(http.StatusOK, gin.H{"cities": jsonCities})
		}
	}
}
