package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/config"
	"plotlyst/server/controllers"
	"plotlyst/server/intelligence"
	"plotlyst/server/middleware"
	"plotlyst/server/tmdb"
	"time"
)

func SetupIntelligenceRoutes(router *gin.Engine, client *mongo.Client, catalog *tmdb.Client, cfg config.Config) {
	local := intelligence.LocalExtractor{}
	var extractor intelligence.Extractor = local
	if cfg.AIProvider == "openai" {
		extractor = intelligence.OpenAIExtractor{Key: cfg.OpenAIAPIKey, Model: cfg.AIModel, MaxTokens: cfg.AIMaxOutputTokens}
	}
	service := &intelligence.Service{Extractor: extractor, Fallback: local, TMDB: catalog, Timeout: cfg.AITimeout, MaxInput: cfg.AIMaxInputChars, DailyBudget: cfg.AIDailyUserBudget, Model: cfg.AIModel, CacheTTL: cfg.AICacheTTL}
	h := controllers.NewIntelligenceController(client, service)
	publicAI := router.Group("/api/ai")
	publicAI.Use(middleware.OptionalAuth())
	publicAI.POST("/discover", middleware.RateLimit(cfg.AIDailyUserBudget, 24*time.Hour), h.Discover("discover"))
	publicAI.POST("/memory-search", middleware.RateLimit(cfg.AIDailyUserBudget, 24*time.Hour), h.Discover("memory"))
	g := router.Group("/api")
	g.Use(middleware.AuthMiddleWare(), middleware.RequireActiveUser(client))
	g.GET("/analytics/cinema-dna", h.CinemaDNA)
	g.GET("/me/ai-privacy", h.GetAIPrivacy)
	g.PUT("/me/ai-privacy", h.PutAIPrivacy)
	g.POST("/movie-night/rooms", h.CreateNight)
	g.POST("/movie-night/rooms/:code/join", h.JoinNight)
	g.GET("/movie-night/rooms/:code", h.GetNight)
	g.PUT("/movie-night/rooms/:code/preferences", h.Preferences)
	g.GET("/movie-night/rooms/:code/ws", h.NightSocket)
}
