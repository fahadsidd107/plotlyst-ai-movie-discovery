package tmdb

import (
	"context"
	"io"
	"net/http"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (fn roundTripFunc) RoundTrip(request *http.Request) (*http.Response, error) { return fn(request) }
func response(status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(body))}
}

func TestClientCoalescesConcurrentRequests(t *testing.T) {
	var calls atomic.Int32
	client := New("https://example.test", "secret", time.Second, time.Minute, 0)
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		time.Sleep(30 * time.Millisecond)
		return response(200, `{"page":1,"results":[]}`), nil
	})
	var wg sync.WaitGroup
	start := make(chan struct{})
	for range 12 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			var page Page[MediaSummary]
			if err := client.Get(context.Background(), "/trending/movie/week", nil, &page); err != nil {
				t.Error(err)
			}
		}()
	}
	close(start)
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("expected one coalesced upstream request, got %d", calls.Load())
	}
}

func TestClientCachesSuccessfulResponses(t *testing.T) {
	var calls atomic.Int32
	client := New("https://example.test", "secret", time.Second, time.Minute, 0)
	client.http.Transport = roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls.Add(1)
		if r.URL.Query().Get("api_key") != "secret" {
			t.Error("missing API key")
		}
		return response(200, `{"id":155,"title":"The Dark Knight"}`), nil
	})
	for range 2 {
		var movie MovieDetail
		if err := client.Get(context.Background(), "/movie/155", nil, &movie); err != nil {
			t.Fatal(err)
		}
		if movie.ID != 155 {
			t.Fatalf("unexpected movie: %+v", movie)
		}
	}
	if calls.Load() != 1 {
		t.Fatalf("expected one upstream call, got %d", calls.Load())
	}
}

func TestClientRetriesServerErrors(t *testing.T) {
	var calls atomic.Int32
	client := New("https://example.test", "secret", time.Second, time.Minute, 2)
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		if calls.Add(1) < 2 {
			return response(503, "try again"), nil
		}
		return response(200, `{"page":1,"results":[]}`), nil
	})
	var page Page[MediaSummary]
	if err := client.Get(context.Background(), "/movie/popular", nil, &page); err != nil {
		t.Fatal(err)
	}
	if calls.Load() != 2 {
		t.Fatalf("expected retry, got %d calls", calls.Load())
	}
}

func TestClientDoesNotRetryNotFound(t *testing.T) {
	var calls atomic.Int32
	client := New("https://example.test", "secret", time.Second, time.Minute, 3)
	client.http.Transport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls.Add(1)
		return response(404, `{"status_message":"missing"}`), nil
	})
	var movie MovieDetail
	if err := client.Get(context.Background(), "/movie/0", nil, &movie); err == nil {
		t.Fatal("expected error")
	}
	if calls.Load() != 1 {
		t.Fatalf("expected no retry, got %d calls", calls.Load())
	}
}
