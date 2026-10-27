package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRoles(t *testing.T) {
	gin.SetMode(gin.TestMode)
	for _, tc := range []struct {
		role string
		want int
	}{{"USER", http.StatusForbidden}, {"MODERATOR", http.StatusNoContent}, {"ADMIN", http.StatusNoContent}} {
		r := gin.New()
		r.Use(func(c *gin.Context) { c.Set("role", tc.role); c.Next() })
		r.GET("/", RequireRoles("MODERATOR", "ADMIN"), func(c *gin.Context) { c.Status(http.StatusNoContent) })
		w := httptest.NewRecorder()
		r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/", nil))
		if w.Code != tc.want {
			t.Fatalf("role %s got %d want %d", tc.role, w.Code, tc.want)
		}
	}
}
