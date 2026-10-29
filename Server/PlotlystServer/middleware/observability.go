package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type metricKey struct {
	Method, Path string
	Status       int
}
type Metrics struct {
	started   time.Time
	mu        sync.RWMutex
	requests  map[metricKey]uint64
	durations map[metricKey]time.Duration
}

func NewMetrics() *Metrics {
	return &Metrics{started: time.Now().UTC(), requests: map[metricKey]uint64{}, durations: map[metricKey]time.Duration{}}
}

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := strings.TrimSpace(c.GetHeader("X-Request-ID"))
		if len(id) < 8 || len(id) > 128 {
			id = randomID()
		}
		c.Set("requestId", id)
		c.Header("X-Request-ID", id)
		c.Next()
	}
}
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("X-Content-Type-Options", "nosniff")
		c.Header("X-Frame-Options", "DENY")
		c.Header("Referrer-Policy", "no-referrer")
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=()")
		if c.Request.TLS != nil || strings.EqualFold(c.GetHeader("X-Forwarded-Proto"), "https") {
			c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		}
		if c.Request.Body != nil {
			c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 2<<20)
		}
		c.Next()
	}
}
func randomID() string {
	b := make([]byte, 12)
	if _, err := rand.Read(b); err != nil {
		return strconv.FormatInt(time.Now().UnixNano(), 36)
	}
	return hex.EncodeToString(b)
}

func (m *Metrics) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		started := time.Now()
		c.Next()
		path := c.FullPath()
		if path == "" {
			path = "unmatched"
		}
		key := metricKey{c.Request.Method, path, c.Writer.Status()}
		elapsed := time.Since(started)
		m.mu.Lock()
		m.requests[key]++
		m.durations[key] += elapsed
		m.mu.Unlock()
		entry := map[string]any{"timestamp": time.Now().UTC().Format(time.RFC3339Nano), "level": "info", "message": "request_completed", "request_id": c.GetString("requestId"), "method": c.Request.Method, "path": path, "status": c.Writer.Status(), "duration_ms": float64(elapsed.Microseconds()) / 1000, "client_ip": c.ClientIP()}
		raw, _ := json.Marshal(entry)
		log.Print(string(raw))
	}
}

func (m *Metrics) Handler(c *gin.Context) {
	m.mu.RLock()
	keys := make([]metricKey, 0, len(m.requests))
	for key := range m.requests {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if keys[i].Path == keys[j].Path {
			return keys[i].Status < keys[j].Status
		}
		return keys[i].Path < keys[j].Path
	})
	var b strings.Builder
	b.WriteString("# HELP plotlyst_uptime_seconds Process uptime in seconds.\n# TYPE plotlyst_uptime_seconds gauge\n")
	fmt.Fprintf(&b, "plotlyst_uptime_seconds %.0f\n", time.Since(m.started).Seconds())
	b.WriteString("# HELP plotlyst_http_requests_total Completed HTTP requests.\n# TYPE plotlyst_http_requests_total counter\n# HELP plotlyst_http_request_duration_seconds_total Cumulative request duration.\n# TYPE plotlyst_http_request_duration_seconds_total counter\n")
	for _, key := range keys {
		labels := fmt.Sprintf("method=%q,path=%q,status=%q", key.Method, key.Path, strconv.Itoa(key.Status))
		fmt.Fprintf(&b, "plotlyst_http_requests_total{%s} %d\n", labels, m.requests[key])
		fmt.Fprintf(&b, "plotlyst_http_request_duration_seconds_total{%s} %.6f\n", labels, m.durations[key].Seconds())
	}
	m.mu.RUnlock()
	c.Data(http.StatusOK, "text/plain; version=0.0.4; charset=utf-8", []byte(b.String()))
}
