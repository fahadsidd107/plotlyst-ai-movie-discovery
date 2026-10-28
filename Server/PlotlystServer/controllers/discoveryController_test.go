package controllers

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func discoveryContext(rawQuery string) *gin.Context {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)
	context.Request = httptest.NewRequest("GET", "/?"+rawQuery, nil)
	return context
}

func TestBuildDiscoverQueryMapsMovieFilters(t *testing.T) {
	query, err := buildDiscoverQuery(discoveryContext("genre=18&year=1999&country=US&language=en&runtime_min=80&runtime_max=140&provider=8&region=de"), "movie")
	if err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{"with_genres": "18", "primary_release_year": "1999", "with_origin_country": "US", "with_original_language": "en", "with_runtime.gte": "80", "with_runtime.lte": "140", "with_watch_providers": "8", "watch_region": "DE"}
	for key, value := range expected {
		if query.Get(key) != value {
			t.Errorf("%s=%q, want %q", key, query.Get(key), value)
		}
	}
}

func TestBuildDiscoverQueryMapsTVDecade(t *testing.T) {
	query, err := buildDiscoverQuery(discoveryContext("decade=1980"), "tv")
	if err != nil {
		t.Fatal(err)
	}
	if query.Get("first_air_date.gte") != "1980-01-01" || query.Get("first_air_date.lte") != "1989-12-31" {
		t.Fatalf("unexpected decade query: %v", query)
	}
}

func TestBuildDiscoverQueryRejectsInvalidRuntime(t *testing.T) {
	if _, err := buildDiscoverQuery(discoveryContext("runtime_min=forever"), "movie"); err == nil {
		t.Fatal("expected invalid runtime error")
	}
}

func TestWeightedScoreRewardsConfidence(t *testing.T) {
	popular := weightedScore(8.0, 10000)
	obscure := weightedScore(8.0, 10)
	if popular <= obscure {
		t.Fatalf("popular score %f should exceed obscure score %f", popular, obscure)
	}
}

func TestDiscoverMonetizationMapping(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest("GET", "/?monetization=subscription", nil)
	q, err := buildDiscoverQuery(c, "movie")
	if err != nil {
		t.Fatal(err)
	}
	if got := q.Get("with_watch_monetization_types"); got != "flatrate" {
		t.Fatalf("got %q", got)
	}
}
