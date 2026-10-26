package routes

import (
	"github.com/gin-gonic/gin"
	"plotlyst/server/controllers"
)

func SetupCatalogRoutes(router *gin.Engine, catalog *controllers.CatalogController) {
	api := router.Group("/api/catalog")
	api.GET("/movies", catalog.ListMovies)
	api.GET("/movies/:id", catalog.Movie)
	api.GET("/tv", catalog.ListTV)
	api.GET("/tv/:id", catalog.TV)
	api.GET("/tv/:id/seasons/:season", catalog.Season)
	api.GET("/tv/:id/seasons/:season/episodes/:episode", catalog.Episode)
	api.GET("/people/trending", catalog.TrendingPeople)
	api.GET("/people/born-today", catalog.BornToday)
	api.GET("/people/:id", catalog.Person)
	api.GET("/search", catalog.Search)
}