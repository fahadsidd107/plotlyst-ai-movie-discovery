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
			return
		}
		if token == "" {
			apiresponse.Error(c, http.StatusUnauthorized, "authentication_required", "Authentication is required")
			return
		}
		claims, err := utils.ValidateToken(token)

		if err != nil {
			apiresponse.Error(c, http.StatusUnauthorized, "invalid_token", "The access token is invalid or expired")
			return
		}
		c.Set("userId", claims.UserId)
		c.Set("role", claims.Role)

		c.Next()

	}
}
