package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"plotlyst/server/apiresponse"
	"plotlyst/server/database"
)

func RequireActiveUser(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		uid, ok := c.Get("userId")
		if !ok {
			c.Abort()
			return
		}
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		var user struct {
			Role      string `bson:"role"`
			Suspended bool   `bson:"suspended"`
		}
		if database.OpenCollection("users", client).FindOne(ctx, bson.M{"user_id": uid}).Decode(&user) != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "session_user_not_found", "The session user no longer exists")
			c.Abort()
			return
		}
		if user.Suspended {
			apiresponse.Error(c, http.StatusForbidden, "account_suspended", "This account is suspended")
			c.Abort()
			return
		}
		c.Set("role", user.Role)
		c.Next()
	}
}
