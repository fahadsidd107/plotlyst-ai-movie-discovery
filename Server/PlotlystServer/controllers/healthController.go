package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/apiresponse"
)

func Liveness(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "plotlyst-api"})
}

func Readiness(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := client.Ping(ctx, nil); err != nil {
			apiresponse.Error(c, http.StatusServiceUnavailable, "service_unavailable", "The database is not ready")
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "ready", "service": "plotlyst-api", "database": "connected"})
	}
}
