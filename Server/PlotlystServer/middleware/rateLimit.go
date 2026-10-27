package middleware

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
)

type rateBucket struct {
	count int
	reset time.Time
}

func RateLimit(limit int, window time.Duration) gin.HandlerFunc {
	var mu sync.Mutex
	buckets := map[string]rateBucket{}
	return func(c *gin.Context) {
		key := c.ClientIP() + ":" + c.FullPath()
		if uid, ok := c.Get("userId"); ok {
			key = uid.(string) + ":" + c.FullPath()
		}
		now := time.Now()
		mu.Lock()
		b := buckets[key]
		if now.After(b.reset) {
			b = rateBucket{reset: now.Add(window)}
		}
		if b.count >= limit {
			retry := int(time.Until(b.reset).Seconds()) + 1
			mu.Unlock()
			c.Header("Retry-After", strconv.Itoa(retry))
			apiresponse.Error(c, http.StatusTooManyRequests, "rate_limited", "Too many requests; please try again shortly")
			c.Abort()
			return
		}
		b.count++
		buckets[key] = b
		mu.Unlock()
		c.Next()
	}
}
