package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/utils"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication is required")
			c.Abort()
			return
		}
		if token == "" {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication is required")
			c.Abort()
			return
		}
		claims, err := utils.ValidateToken(token)

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "invalid_token", "The access token is invalid or expired")
			c.Abort()
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()

	}
}

// OptionalAuth enriches public requests when a valid session is present but
// never turns an anonymous or expired session into an authentication error.
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)
		if err == nil && token != "" {
			if claims, validateErr := utils.ValidateToken(token); validateErr == nil {
				c.Set("userId", claims.UserId)
				c.Set("role", claims.Role)
			}
		}
		c.Next()
	}
}
