package tmdb

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/sync/singleflight"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

type APIError struct {
	Status  int
	Message string
}

func (e *APIError) Error() string { return fmt.Sprintf("TMDB returned %d: %s", e.Status, e.Message) }

type cacheItem struct {
	body    []byte
	expires time.Time
}
type Client struct {
	baseURL, apiKey string
	http            *http.Client
	ttl             time.Duration
	retries         int
	mu              sync.RWMutex
	cache           map[string]cacheItem
	inflight        singleflight.Group
}

func New(baseURL, apiKey string, timeout, ttl time.Duration, retries int) *Client {
	return &Client{baseURL: strings.TrimRight(baseURL, "/"), apiKey: apiKey, http: &http.Client{Timeout: timeout}, ttl: ttl, retries: retries, cache: make(map[string]cacheItem)}
}

func (c *Client) Get(ctx context.Context, path string, query url.Values, destination any) error {
	if query == nil {
		query = url.Values{}
	} else {
		query = cloneValues(query)
	}
	query.Set("api_key", c.apiKey)
	if !query.Has("language") {
		query.Set("language", "en-US")
	}
	requestURL := c.baseURL + path + "?" + query.Encode()
	if body, ok := c.cached(requestURL); ok {
		return json.Unmarshal(body, destination)
	}
	value, err, _ := c.inflight.Do(requestURL, func() (any, error) {
		if body, ok := c.cached(requestURL); ok {
			return body, nil
		}
		var lastErr error
		for attempt := 0; attempt <= c.retries; attempt++ {
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, requestURL, nil)
			if err != nil {
				return nil, err
			}
			req.Header.Set("Accept", "application/json")
			resp, err := c.http.Do(req)
			if err != nil {
				lastErr = err
				if !sleep(ctx, backoff(attempt)) {
					return nil, ctx.Err()
				}
				continue
			}
			body, readErr := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
			resp.Body.Close()
			if readErr != nil {
				return nil, readErr
			}
			if resp.StatusCode >= 200 && resp.StatusCode < 300 {
				c.store(requestURL, body)
				return body, nil
			}
			message := strings.TrimSpace(string(body))
			var apiMessage struct {
				StatusMessage string `json:"status_message"`
			}
			if json.Unmarshal(body, &apiMessage) == nil && apiMessage.StatusMessage != "" {
				message = apiMessage.StatusMessage
			}
			lastErr = &APIError{Status: resp.StatusCode, Message: message}
			if resp.StatusCode != http.StatusTooManyRequests && resp.StatusCode < 500 {
				return nil, lastErr
			}
			delay := backoff(attempt)
			if seconds, err := strconv.Atoi(resp.Header.Get("Retry-After")); err == nil && seconds > 0 {
				delay = time.Duration(seconds) * time.Second
			}
			if !sleep(ctx, delay) {
				return nil, ctx.Err()
			}
		}
		return nil, lastErr
	})
	if err != nil {
		return err
	}
	return json.Unmarshal(value.([]byte), destination)
}

func (c *Client) cached(key string) ([]byte, bool) {
	c.mu.RLock()
	item, ok := c.cache[key]
	c.mu.RUnlock()
	if !ok || time.Now().After(item.expires) {
		if ok {
			c.mu.Lock()
			delete(c.cache, key)
			c.mu.Unlock()
		}
		return nil, false
	}
	return append([]byte(nil), item.body...), true
}
func (c *Client) store(key string, body []byte) {
	c.mu.Lock()
	if len(c.cache) >= 500 {
		now := time.Now()
		for candidate, item := range c.cache {
			if now.After(item.expires) {
				delete(c.cache, candidate)
			}
		}
		if len(c.cache) >= 500 {
			for candidate := range c.cache {
				delete(c.cache, candidate)
				break
			}
		}
	}
	c.cache[key] = cacheItem{body: append([]byte(nil), body...), expires: time.Now().Add(c.ttl)}
	c.mu.Unlock()
}
func cloneValues(source url.Values) url.Values {
	copy := url.Values{}
	for key, values := range source {
		copy[key] = append([]string(nil), values...)
	}
	return copy
}
func backoff(attempt int) time.Duration {
	return time.Duration(1<<min(attempt, 3)) * 200 * time.Millisecond
}
func sleep(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
