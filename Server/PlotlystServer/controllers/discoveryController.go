package controllers

import (
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/tmdb"
)

type DiscoveryController struct{ catalog *CatalogController }

func NewDiscoveryController(catalog *CatalogController) *DiscoveryController {
	return &DiscoveryController{catalog: catalog}
}

type homeSection struct {
	Key       string              `json:"key"`
	Title     string              `json:"title"`
	MediaType string              `json:"media_type"`
	Results   []tmdb.MediaSummary `json:"results"`
	Error     string              `json:"error,omitempty"`
}
type chartItem struct {
	tmdb.MediaSummary
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
}

func (h *DiscoveryController) Home(c *gin.Context) {
	definitions := []struct{ key, title, media, path string }{
		{"now_playing", "Now playing in cinemas", "movie", "/movie/now_playing"},
		{"trending_movies", "Movies moving culture", "movie", "/trending/movie/week"},
		{"trending_tv", "Series everyone is discussing", "tv", "/trending/tv/week"},
		{"fan_favourites", "Fan favourites", "movie", "/movie/top_rated"},
	}
	sections := make([]homeSection, len(definitions))
	var wg sync.WaitGroup
	for i, definition := range definitions {
		i, definition := i, definition
		wg.Add(1)
		go func() {
			defer wg.Done()
			var page tmdb.Page[tmdb.MediaSummary]
			err := h.catalog.client.Get(c.Request.Context(), definition.path, url.Values{"page": {"1"}, "region": {c.DefaultQuery("region", "DE")}}, &page)
			sections[i] = homeSection{Key: definition.key, Title: definition.title, MediaType: definition.media, Results: page.Results}
			if err != nil {
				sections[i].Error = "This reel is temporarily unavailable."
			}
		}()
	}
	wg.Wait()
	if h.catalog.mode(c) == "STRICT" {
		for i := range sections {
			protectMedia(sections[i].Results)
		}
	}
	h.catalog.respond(c, gin.H{"generated_at": time.Now().UTC(), "sections": sections})
}

func (h *DiscoveryController) HomeSection(c *gin.Context) {
	section := c.Param("section")
	definitions := map[string]struct{ title, media, path string }{
		"now-playing": {"Now playing in cinemas", "movie", "/movie/now_playing"}, "trending-movies": {"Movies moving culture", "movie", "/trending/movie/week"}, "trending-tv": {"Series everyone is discussing", "tv", "/trending/tv/week"}, "fan-favourites": {"Fan favourites", "movie", "/movie/top_rated"}, "airing-today": {"Airing today", "tv", "/tv/airing_today"},
	}
	definition, ok := definitions[section]
	if !ok {
		apiresponse.Error(c, http.StatusNotFound, "discovery_section_not_found", "The requested discovery section does not exist")
		return
	}
	var page tmdb.Page[tmdb.MediaSummary]
	if !h.catalog.fetch(c, definition.path, listQuery(c), &page) {
		return
	}
	if h.catalog.mode(c) == "STRICT" {
		protectMedia(page.Results)
	}
	h.catalog.respond(c, homeSection{Key: section, Title: definition.title, MediaType: definition.media, Results: page.Results})
}

func (h *DiscoveryController) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len([]rune(query)) < 2 {
		apiresponse.Error(c, http.StatusBadRequest, "query_too_short", "Search requires at least two characters")
		return
	}
	var page tmdb.Page[tmdb.MediaSummary]
	q := url.Values{"query": {query}, "page": {validPage(c.Query("page"))}, "include_adult": {"false"}}
	if !h.catalog.fetch(c, "/search/multi", q, &page) {
		return
	}
	filtered := make([]tmdb.MediaSummary, 0, len(page.Results))
	for _, item := range page.Results {
		if item.MediaType == "movie" || item.MediaType == "tv" || item.MediaType == "person" {
			filtered = append(filtered, item)
		}
	}
	page.Results = filtered
	if h.catalog.mode(c) == "STRICT" {
		protectMedia(page.Results)
	}
	h.catalog.respond(c, page)
}

func (h *DiscoveryController) DiscoverMovies(c *gin.Context) { h.discover(c, "movie") }
func (h *DiscoveryController) DiscoverTV(c *gin.Context)     { h.discover(c, "tv") }

func (h *DiscoveryController) FilterOptions(c *gin.Context) {
	media := c.Param("media")
	if media != "movie" && media != "tv" {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_media_type", "Filter media must be movie or tv")
		return
	}
	var genres struct {
		Genres []tmdb.Genre `json:"genres"`
	}
	var providers struct {
		Results []tmdb.Provider `json:"results"`
	}
	var countries []struct {
		ISO31661    string `json:"iso_3166_1"`
		EnglishName string `json:"english_name"`
	}
	var languages []struct {
		ISO6391     string `json:"iso_639_1"`
		EnglishName string `json:"english_name"`
		Name        string `json:"name"`
	}
	type fetchResult struct{ err error }
	results := make(chan fetchResult, 4)
	ctx := c.Request.Context()
	go func() { results <- fetchResult{h.catalog.client.Get(ctx, "/genre/"+media+"/list", nil, &genres)} }()
	go func() {
		results <- fetchResult{h.catalog.client.Get(ctx, "/watch/providers/"+media, url.Values{"watch_region": {c.DefaultQuery("region", "DE")}}, &providers)}
	}()
	go func() { results <- fetchResult{h.catalog.client.Get(ctx, "/configuration/countries", nil, &countries)} }()
	go func() { results <- fetchResult{h.catalog.client.Get(ctx, "/configuration/languages", nil, &languages)} }()
	for range 4 {
		if entry := <-results; entry.err != nil {
			apiresponse.Error(c, http.StatusBadGateway, "filter_options_unavailable", "Discovery filter options are temporarily unavailable")
			return
		}
	}
	runtimeMax := 360
	if media == "tv" {
		runtimeMax = 240
	}
	h.catalog.respond(c, gin.H{"media_type": media, "genres": genres.Genres, "providers": providers.Results, "countries": countries, "languages": languages, "runtime": gin.H{"min": 0, "max": runtimeMax}, "year": gin.H{"min": 1880, "max": time.Now().Year() + 5}})
}
func (h *DiscoveryController) discover(c *gin.Context, media string) {
	q, err := buildDiscoverQuery(c, media)
	if err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_filter", err.Error())
		return
	}
	var page tmdb.Page[tmdb.MediaSummary]
	if !h.catalog.fetch(c, "/discover/"+media, q, &page) {
		return
	}
	if h.catalog.mode(c) == "STRICT" {
		protectMedia(page.Results)
	}
	h.catalog.respond(c, page)
}

func (h *DiscoveryController) Chart(c *gin.Context) {
	media := c.Param("media")
	if media != "movie" && media != "tv" {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_media_type", "Chart media must be movie or tv")
		return
	}
	size, err := strconv.Atoi(c.Param("size"))
	if err != nil || (size != 100 && size != 250 && size != 2000) {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_chart_size", "Chart size must be 100, 250, or 2000")
		return
	}
	pageNum, _ := strconv.Atoi(validPage(c.Query("page")))
	maxPages := (size + 19) / 20
	if pageNum > maxPages {
		apiresponse.Error(c, http.StatusBadRequest, "chart_page_out_of_range", fmt.Sprintf("This chart has %d pages", maxPages))
		return
	}
	q := url.Values{"page": {strconv.Itoa(pageNum)}, "sort_by": {"vote_average.desc"}, "vote_count.gte": {"500"}, "include_adult": {"false"}}
	var page tmdb.Page[tmdb.MediaSummary]
	if !h.catalog.fetch(c, "/discover/"+media, q, &page) {
		return
	}
	if h.catalog.mode(c) == "STRICT" {
		protectMedia(page.Results)
	}
	items := make([]chartItem, 0, len(page.Results))
	for _, item := range page.Results {
		items = append(items, chartItem{MediaSummary: item, Score: weightedScore(item.VoteAverage, item.VoteCount)})
	}
	sort.SliceStable(items, func(i, j int) bool {
		if items[i].Score != items[j].Score {
			return items[i].Score > items[j].Score
		}
		if items[i].VoteCount != items[j].VoteCount {
			return items[i].VoteCount > items[j].VoteCount
		}
		return items[i].ID < items[j].ID
	})
	for i := range items {
		items[i].Rank = (pageNum-1)*20 + i + 1
	}
	h.catalog.respond(c, gin.H{"media_type": media, "size": size, "page": pageNum, "total_pages": maxPages, "items": items, "methodology": gin.H{"formula": "(v/(v+m))*R + (m/(v+m))*C", "minimum_votes": 500, "global_mean": 6.5, "tie_break": "score desc, vote_count desc, id asc", "source": "TMDB vote_average and vote_count", "snapshot": "live cached provider response"}})
}

func (h *DiscoveryController) Collection(c *gin.Context) {
	if !validNumericID(c) {
		return
	}
	var out struct {
		ID           int                 `json:"id"`
		Name         string              `json:"name"`
		Overview     string              `json:"overview"`
		PosterPath   string              `json:"poster_path"`
		BackdropPath string              `json:"backdrop_path"`
		Parts        []tmdb.MediaSummary `json:"parts"`
	}
	if !h.catalog.fetch(c, "/collection/"+c.Param("id"), nil, &out) {
		return
	}
	if h.catalog.mode(c) == "STRICT" {
		out.Overview = ""
		protectMedia(out.Parts)
	}
	h.catalog.respond(c, gin.H{"kind": "official_franchise", "ownership": "TMDB collection", "collection": out})
}

func (h *DiscoveryController) CollectionKinds(c *gin.Context) {
	h.catalog.respond(c, []gin.H{{"kind": "official_franchise", "owner": "TMDB", "description": "Canonical connected-film collections supplied by TMDB."}, {"kind": "editorial_collection", "owner": "Plotlyst editors", "description": "Curated themes; stored separately from official franchises."}, {"kind": "user_list", "owner": "Plotlyst user", "description": "User-owned lists with independent permissions; persistence arrives with Step 4."}})
}

func buildDiscoverQuery(c *gin.Context, media string) (url.Values, error) {
	q := url.Values{"page": {validPage(c.Query("page"))}, "include_adult": {"false"}, "sort_by": {allowedSort(c.DefaultQuery("sort", "popularity.desc"))}}
	copyIf := func(source, target string) {
		if value := strings.TrimSpace(c.Query(source)); value != "" {
			q.Set(target, value)
		}
	}
	copyIf("genre", "with_genres")
	copyIf("country", "with_origin_country")
	copyIf("language", "with_original_language")
	copyIf("runtime_min", "with_runtime.gte")
	copyIf("runtime_max", "with_runtime.lte")
	copyIf("provider", "with_watch_providers")
	if monetization := c.Query("monetization"); monetization != "" {
		allowed := map[string]string{"subscription": "flatrate", "rental": "rent", "purchase": "buy", "free": "free|ads"}
		mapped, ok := allowed[monetization]
		if !ok {
			return nil, fmt.Errorf("monetization must be subscription, rental, purchase, or free")
		}
		q.Set("with_watch_monetization_types", mapped)
	}
	region := strings.ToUpper(c.DefaultQuery("region", "DE"))
	if len(region) != 2 {
		return nil, fmt.Errorf("region must be a two-letter country code")
	}
	q.Set("region", region)
	q.Set("watch_region", region)
	if year := c.Query("year"); year != "" {
		value, err := strconv.Atoi(year)
		if err != nil || value < 1880 || value > time.Now().Year()+5 {
			return nil, fmt.Errorf("year must be between 1880 and %d", time.Now().Year()+5)
		}
		if media == "movie" {
			q.Set("primary_release_year", year)
		} else {
			q.Set("first_air_date_year", year)
		}
	}
	if decade := c.Query("decade"); decade != "" {
		value, err := strconv.Atoi(decade)
		if err != nil || value < 1880 || value%10 != 0 {
			return nil, fmt.Errorf("decade must be a year ending in 0")
		}
		prefix := "primary_release_date"
		if media == "tv" {
			prefix = "first_air_date"
		}
		q.Set(prefix+".gte", fmt.Sprintf("%d-01-01", value))
		q.Set(prefix+".lte", fmt.Sprintf("%d-12-31", value+9))
	}
	for _, key := range []string{"with_runtime.gte", "with_runtime.lte"} {
		if raw := q.Get(key); raw != "" {
			n, err := strconv.Atoi(raw)
			if err != nil || n < 0 || n > 600 {
				return nil, fmt.Errorf("runtime must be between 0 and 600 minutes")
			}
		}
	}
	return q, nil
}
func validPage(raw string) string {
	n, err := strconv.Atoi(raw)
	if err != nil || n < 1 || n > 500 {
		return "1"
	}
	return strconv.Itoa(n)
}
func allowedSort(raw string) string {
	allowed := map[string]bool{"popularity.desc": true, "vote_average.desc": true, "vote_count.desc": true, "primary_release_date.desc": true, "first_air_date.desc": true}
	if allowed[raw] {
		return raw
	}
	return "popularity.desc"
}
func weightedScore(rating float64, votes int) float64 {
	const minimum = 500.0
	const mean = 6.5
	v := float64(votes)
	return (v/(v+minimum))*rating + (minimum/(v+minimum))*mean
}
