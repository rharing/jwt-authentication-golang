package main

import (
	"jwt-authentication-golang/config"
	"jwt-authentication-golang/controllers"
	"jwt-authentication-golang/database"
	"jwt-authentication-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	dbUsername := config.Config("db_username")
	dbPassword := config.Config("db_password")
	connectUrl := dbUsername + ":" + dbPassword + "@tcp(localhost:3306)/jwt_demo?parseTime=true"
	database.Connect(connectUrl)
	database.Migrate()

	// Initialize Router
	router := initRouter()
	router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()

	api := router.Group("/api")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
		movies := api.Group("/movies").Use(middlewares.Auth())
		{
			movies.GET("/cities", controllers.GetCities)
		}
		movie := api.Group("/movie/:movieId").Use(middlewares.Auth())
		{
			movie.GET("/", controllers.GetMovie)
		}
		cityMovies := api.Group("/movies/:city").Use(middlewares.Auth())
		{
			cityMovies.GET("/", controllers.GetCity)
		}
	}
	return router
}
