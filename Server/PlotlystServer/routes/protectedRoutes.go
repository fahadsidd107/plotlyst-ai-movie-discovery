package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	controller "plotlyst/server/controllers"
	"plotlyst/server/middleware"
)

func SetupProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	protected := router.Group("")
	protected.Use(middleware.AuthMiddleWare())

	protected.GET("/movie/:imdb_id", controller.GetMovie(client))
	protected.POST("/addmovie", controller.AddMovie(client))
	protected.GET("/recommendedmovies", controller.GetRecommendedMovies(client))
	protected.PATCH("/updatereview/:imdb_id", controller.AdminReviewUpdate(client))
	protected.POST("/logout", controller.LogoutHandler(client))
	protected.POST("/api/auth/logout", controller.LogoutHandler(client))
}
