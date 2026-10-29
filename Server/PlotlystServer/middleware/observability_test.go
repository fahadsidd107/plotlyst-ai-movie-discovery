package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequestIDAndMetrics(t *testing.T) {
	gin.SetMode(gin.TestMode)
	metrics := NewMetrics()
	r := gin.New()
	r.Use(RequestID(), metrics.Middleware())
	r.GET("/probe", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })
	r.GET("/metrics", metrics.Handler)
	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/probe", nil)
	req.Header.Set("X-Request-ID", "test-request-123")
	r.ServeHTTP(w, req)
	if w.Header().Get("X-Request-ID") != "test-request-123" {
		t.Fatalf("request id not propagated: %q", w.Header().Get("X-Request-ID"))
	}
	w = httptest.NewRecorder()
	r.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/metrics", nil))
	body := w.Body.String()
	if !strings.Contains(body, `plotlyst_http_requests_total{method="GET",path="/probe",status="200"} 1`) {
		t.Fatalf("metrics missing request: %s", body)
	}
}
