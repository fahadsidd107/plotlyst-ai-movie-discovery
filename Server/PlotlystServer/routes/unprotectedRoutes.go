package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	controller "plotlyst/server/controllers"
	"plotlyst/server/middleware"
	"time"
)

func SetupUnProtectedRoutes(router *gin.Engine, client *mongo.Client) {

	router.GET("/movies", controller.GetMovies(client))
	router.POST("/register", middleware.RateLimit(8, time.Hour), controller.RegisterUser(client))
	router.POST("/login", middleware.RateLimit(12, time.Minute), controller.LoginUser(client))
	router.GET("/genres", controller.GetGenres(client))
	router.POST("/refresh", controller.RefreshTokenHandler(client))
	router.POST("/api/auth/register", middleware.RateLimit(8, time.Hour), controller.RegisterUser(client))
	router.POST("/api/auth/login", middleware.RateLimit(12, time.Minute), controller.LoginUser(client))
	router.POST("/api/auth/refresh", controller.RefreshTokenHandler(client))
}
