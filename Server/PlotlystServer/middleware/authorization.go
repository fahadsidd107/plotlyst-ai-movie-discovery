package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
)

func RequireRoles(roles ...string) gin.HandlerFunc {
	allowed := map[string]bool{}
	for _, role := range roles {
		allowed[strings.ToUpper(role)] = true
	}
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists || !allowed[strings.ToUpper(role.(string))] {
			apiresponse.Error(c, http.StatusForbidden, "insufficient_role", "Your account cannot perform this action")
			c.Abort()
			return
		}
		c.Next()
	}
}
