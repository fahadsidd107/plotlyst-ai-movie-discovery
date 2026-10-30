package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/config"
	"plotlyst/server/controllers"
	"plotlyst/server/database"
	"plotlyst/server/debate"
	"plotlyst/server/middleware"
)

func SetupDebateRoutes(router *gin.Engine, client *mongo.Client, cfg config.Config) {
	local := debate.DeterministicAnalyzer{}
	var ai debate.Analyzer
	if cfg.DebateAIEnabled {
		ai = debate.OpenAIAnalyzer{Key: cfg.OpenAIAPIKey, Model: cfg.AIModel, MaxTokens: cfg.AIMaxOutputTokens, Fallback: local}
	}
	service := &debate.Service{Reviews: database.OpenCollection("reviews", client), Ratings: database.OpenCollection("ratings", client), Cache: database.OpenCollection("debate_summaries", client), AIPrivacy: database.OpenCollection("ai_privacy", client), Local: local, AI: ai, Timeout: cfg.AITimeout}
	h := controllers.NewDebateController(service)
	group := router.Group("/api/debate")
	group.Use(middleware.OptionalAuth())
	group.GET("/movies/:id", h.Movie)
}
