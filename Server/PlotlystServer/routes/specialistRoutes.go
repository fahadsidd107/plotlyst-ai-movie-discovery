package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
	"plotlyst/server/config"
	"plotlyst/server/controllers"
	"plotlyst/server/middleware"
	"plotlyst/server/models"
	"plotlyst/server/specialist"
	"plotlyst/server/tmdb"
)

func SetupSpecialistRoutes(router *gin.Engine, client *mongo.Client, tmdbClient *tmdb.Client, cfg config.Config) {
	var boxOffice specialist.BoxOfficeProvider
	if cfg.OpusDataAPIKey != "" && cfg.OpusDataBaseURL != "" {
		boxOffice = specialist.OpusDataProvider{BaseURL: cfg.OpusDataBaseURL, APIKey: cfg.OpusDataAPIKey, Client: &http.Client{Timeout: cfg.OpusDataTimeout}}
	} else if cfg.BoxOfficeDemoEnabled {
		boxOffice = specialist.DemoBoxOffice{}
	}
	c := controllers.NewSpecialistController(client, specialist.TMDBAvailability{Client: tmdbClient}, boxOffice)
	g := router.Group("/api/specialist")
	g.GET("/providers", c.Providers)
	g.GET("/availability/:mediaType/:id", c.Availability)
	g.GET("/awards/organizations", c.AwardOrganizations)
	g.GET("/awards/:organization/:year", c.AwardCeremony)
	g.GET("/box-office", c.BoxOffice)
	g.GET("/editorial", c.Editorial)
	admin := router.Group("/api/admin/imports")
	admin.Use(middleware.AuthMiddleWare(), middleware.RequireActiveUser(client), middleware.RequireRoles(models.RoleAdmin))
	admin.POST("/awards", c.Import("awards"))
	admin.POST("/box-office", c.Import("box-office"))
	admin.POST("/editorial", c.Import("editorial"))
}
