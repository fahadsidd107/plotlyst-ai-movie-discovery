package routes

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/controllers"
	"plotlyst/server/database"
	"plotlyst/server/middleware"
	"plotlyst/server/storydna"
	"plotlyst/server/tmdb"
)

func SetupStoryDNARoutes(router *gin.Engine, client *mongo.Client, catalogue *tmdb.Client) {
	service := &storydna.Service{TMDB: catalogue, Collection: database.OpenCollection("story_dna_profiles", client)}
	h := controllers.NewStoryDNAController(service)
	g := router.Group("/api/story-dna")
	g.Use(middleware.OptionalAuth())
	g.GET("/movies/:id", h.Movie)
	g.POST("/discover", middleware.RateLimit(40, 24*time.Hour), h.Discover)
}
