package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/controllers"
	"plotlyst/server/middleware"
	"plotlyst/server/models"
	"time"
)

func SetupCommunityRoutes(router *gin.Engine, client *mongo.Client) {
	cc := controllers.NewCommunityController(client)
	router.GET("/api/profiles/:username", cc.GetPublicProfile)
	router.GET("/api/community/:mediaType/:mediaId/ratings", middleware.AuthMiddleWare(), cc.RatingSummary)
	router.GET("/api/community/:mediaType/:mediaId/reviews", middleware.OptionalAuth(), cc.ListReviews)
	router.GET("/api/community/reviews/:reviewId/reveal", middleware.OptionalAuth(), cc.RevealReview)

	auth := router.Group("/api")
	auth.Use(middleware.AuthMiddleWare())
	auth.Use(middleware.RequireActiveUser(client))
	auth.GET("/auth/session", cc.Session)
	auth.GET("/me/profile", cc.GetMyProfile)
	auth.PUT("/me/profile", cc.UpsertMyProfile)
	auth.PUT("/me/settings", cc.UpsertMyProfile)
	auth.GET("/me/spoiler-preferences", cc.GetSpoilerPreference)
	auth.PUT("/me/spoiler-preferences", cc.PutSpoilerPreference)
	auth.DELETE("/me/account", cc.DeleteAccount)
	auth.PUT("/community/:mediaType/:mediaId/rating", cc.PutRating)
	auth.DELETE("/community/:mediaType/:mediaId/rating", cc.DeleteRating)
	auth.POST("/community/:mediaType/:mediaId/reviews", middleware.RateLimit(8, time.Minute), cc.CreateReview)
	auth.PATCH("/reviews/:reviewId", cc.UpdateReview)
	auth.DELETE("/reviews/:reviewId", cc.DeleteReview)
	auth.POST("/reviews/:reviewId/replies", middleware.RateLimit(12, time.Minute), cc.CreateReply)
	auth.PUT("/reviews/:reviewId/helpful", cc.ToggleHelpful)
	auth.DELETE("/reviews/:reviewId/helpful", cc.RemoveHelpful)
	auth.POST("/reviews/:reviewId/reports", middleware.RateLimit(6, time.Hour), cc.ReportReview)
	auth.GET("/me/watchlist", cc.ListMine("watchlist_items"))
	auth.PUT("/me/watchlist/:mediaType/:mediaId", cc.PutLibrary("watchlist_items"))
	auth.DELETE("/me/watchlist/:mediaType/:mediaId", cc.DeleteLibrary("watchlist_items"))
	auth.GET("/me/history", cc.ListMine("watch_history"))
	auth.PUT("/me/history/:mediaType/:mediaId", cc.PutLibrary("watch_history"))
	auth.DELETE("/me/history/:mediaType/:mediaId", cc.DeleteLibrary("watch_history"))
	auth.GET("/me/playlists", cc.ListPlaylists)
	auth.POST("/me/playlists", cc.CreatePlaylist)
	auth.GET("/me/playlists/:playlistId", cc.GetPlaylist)
	auth.PATCH("/me/playlists/:playlistId", cc.UpdatePlaylist)
	auth.DELETE("/me/playlists/:playlistId", cc.DeletePlaylist)
	auth.POST("/me/playlists/:playlistId/items", cc.AddPlaylistItem)
	auth.DELETE("/me/playlists/:playlistId/items/:itemId", cc.DeletePlaylistItem)
	auth.PUT("/me/tv/:showId/seasons/:season/episodes/:episode/progress", cc.PutProgress)
	auth.GET("/me/tv/:showId/progress", cc.GetProgress)

	moderation := auth.Group("/moderation")
	moderation.Use(middleware.RequireRoles(models.RoleModerator, models.RoleAdmin))
	moderation.GET("/reports", cc.ModerationReports)
	moderation.PATCH("/reviews/:reviewId", cc.ModerateReview)
	moderation.POST("/users/:userId/suspend", middleware.RequireRoles(models.RoleAdmin), cc.SuspendUser)
}
