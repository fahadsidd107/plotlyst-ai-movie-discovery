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
	discovery := controllers.NewDiscoveryController(catalog)
	d := router.Group("/api/discovery")
	d.GET("/home", discovery.Home)
	d.GET("/home/:section", discovery.HomeSection)
	d.GET("/search", discovery.Search)
	d.GET("/movies", discovery.DiscoverMovies)
	d.GET("/tv", discovery.DiscoverTV)
	d.GET("/filters/:media", discovery.FilterOptions)
	d.GET("/charts/:media/:size", discovery.Chart)
	d.GET("/franchises/:id", discovery.Collection)
	d.GET("/collection-kinds", discovery.CollectionKinds)
}
