package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"plotlyst/server/apiresponse"
	"plotlyst/server/spoiler"
	"plotlyst/server/tmdb"
)

type CatalogController struct {
	client       *tmdb.Client
	imageBaseURL string
	spoilers     *spoiler.Service
}

func NewCatalogController(client *tmdb.Client, imageBaseURL string, spoilers *spoiler.Service) *CatalogController {
	return &CatalogController{client: client, imageBaseURL: strings.TrimRight(imageBaseURL, "/"), spoilers: spoilers}
}

func (h *CatalogController) ListMovies(c *gin.Context) {
	var out tmdb.Page[tmdb.MovieSummary]
	q := listQuery(c)
	path := "/movie/popular"
	if c.Query("collection") == "trending" {
		path = "/trending/movie/week"
	}
	if !h.fetch(c, path, q, &out) {
		return
	}
	h.respond(c, out)
}
func (h *CatalogController) ListTV(c *gin.Context) {
	var out tmdb.Page[tmdb.TVSummary]
	q := listQuery(c)
	path := "/tv/popular"
	if c.Query("collection") == "trending" {
		path = "/trending/tv/week"
	}
	if !h.fetch(c, path, q, &out) {
		return
	}
	h.respond(c, out)
}
func (h *CatalogController) Movie(c *gin.Context) {
	var out tmdb.MovieDetail
	if !validNumericID(c) {
		return
	}
	q := url.Values{"append_to_response": {"credits,release_dates,recommendations,watch/providers"}}
	if !h.fetch(c, "/movie/"+c.Param("id"), q, &out) {
		return
	}
	if h.mode(c) == spoiler.ModeStrict && c.Query("reveal_spoilers") != "true" {
		out.Overview = ""
		out.Tagline = ""
		out.SpoilerHidden = true
		protectMedia(out.Recommendations.Results)
	}
	h.respond(c, out)
}
func (h *CatalogController) TV(c *gin.Context) {
	var out tmdb.TVDetail
	if !validNumericID(c) {
		return
	}
	q := url.Values{"append_to_response": {"credits,recommendations,watch/providers,content_ratings"}}
	if !h.fetch(c, "/tv/"+c.Param("id"), q, &out) {
		return
	}
	if h.mode(c) == spoiler.ModeStrict && c.Query("reveal_spoilers") != "true" {
		out.Overview = ""
		out.Tagline = ""
		out.SpoilerHidden = true
		protectMedia(out.Recommendations.Results)
	}
	h.respond(c, out)
}
func (h *CatalogController) Season(c *gin.Context) {
	var out tmdb.SeasonDetail
	if !validNumericID(c) || !validNumber(c, "season") {
		return
	}
	q := url.Values{"append_to_response": {"credits"}}
	if !h.fetch(c, fmt.Sprintf("/tv/%s/season/%s", c.Param("id"), c.Param("season")), q, &out) {
		return
	}
	if c.Query("reveal_spoilers") != "true" {
		show, _ := strconv.Atoi(c.Param("id"))
		season, _ := strconv.Atoi(c.Param("season"))
		mode := h.mode(c)
		for i := range out.Episodes {
			if !h.spoilers.EpisodeVisible(c.Request.Context(), c.GetString("userId"), mode, show, season, out.Episodes[i].EpisodeNumber) {
				out.Episodes[i].Overview = ""
				out.Episodes[i].SpoilerHidden = true
			}
		}
	}
	h.respond(c, out)
}
func (h *CatalogController) Episode(c *gin.Context) {
	var out tmdb.EpisodeDetail
	if !validNumericID(c) || !validNumber(c, "season") || !validNumber(c, "episode") {
		return
	}
	q := url.Values{"append_to_response": {"credits"}}
	if !h.fetch(c, fmt.Sprintf("/tv/%s/season/%s/episode/%s", c.Param("id"), c.Param("season"), c.Param("episode")), q, &out) {
		return
	}
	if c.Query("reveal_spoilers") != "true" {
		show, _ := strconv.Atoi(c.Param("id"))
		season, _ := strconv.Atoi(c.Param("season"))
		episode, _ := strconv.Atoi(c.Param("episode"))
		if !h.spoilers.EpisodeVisible(c.Request.Context(), c.GetString("userId"), h.mode(c), show, season, episode) {
			out.Overview = ""
			out.SpoilerHidden = true
		}
	}
	h.respond(c, out)
}

func (h *CatalogController) mode(c *gin.Context) string {
	if h.spoilers == nil {
		return spoiler.ModeStandard
	}
	return h.spoilers.Get(c.Request.Context(), c.GetString("userId")).Mode
}
func protectMedia(items []tmdb.MediaSummary) {
	for i := range items {
		items[i].Overview = ""
		items[i].SpoilerHidden = true
	}
}
func (h *CatalogController) TrendingPeople(c *gin.Context) {
	var out tmdb.Page[tmdb.PersonSummary]
	if !h.fetch(c, "/trending/person/week", listQuery(c), &out) {
		return
	}
	h.respond(c, out)
}
func (h *CatalogController) Person(c *gin.Context) {
	var out tmdb.PersonDetail
	if !validNumericID(c) {
		return
	}
	q := url.Values{"append_to_response": {"combined_credits,external_ids"}}
	if !h.fetch(c, "/person/"+c.Param("id"), q, &out) {
		return
	}
	h.respond(c, out)
}

func (h *CatalogController) Search(c *gin.Context) {
	var out tmdb.Page[interface{}]
	query := c.Query("query")
	if query == "" {
		apiresponse.Error(c, http.StatusBadRequest, "query_required", "Search query is required")
		return
	}
	q := url.Values{"query": {query}}
	if !h.fetch(c, "/search/multi", q, &out) {
		return
	}
	h.respond(c, out)
}

func (h *CatalogController) BornToday(c *gin.Context) {
	location := time.UTC
	if tz := c.Query("timezone"); tz != "" {
		parsed, err := time.LoadLocation(tz)
		if err != nil {
			apiresponse.Error(c, http.StatusBadRequest, "invalid_timezone", "Use an IANA timezone such as Europe/Berlin")
			return
		}
		location = parsed
	}
	today := time.Now().In(location).Format("01-02")
	var people []tmdb.PersonSummary
	for page := 1; page <= 3; page++ {
		var result tmdb.Page[tmdb.PersonSummary]
		if !h.fetch(c, "/person/popular", url.Values{"page": {strconv.Itoa(page)}}, &result) {
			return
		}
		people = append(people, result.Results...)
	}
	matched := make([]tmdb.PersonDetail, 0)
	var mu sync.Mutex
	var wg sync.WaitGroup
	slots := make(chan struct{}, 6)
	for _, person := range people {
		person := person
		wg.Add(1)
		go func() {
			defer wg.Done()
			slots <- struct{}{}
			defer func() { <-slots }()
			var detail tmdb.PersonDetail
			if h.client.Get(c.Request.Context(), fmt.Sprintf("/person/%d", person.ID), url.Values{"append_to_response": {"combined_credits,external_ids"}}, &detail) == nil && len(detail.Birthday) >= 10 && detail.Birthday[5:] == today {
				mu.Lock()
				matched = append(matched, detail)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	h.respond(c, gin.H{"date": time.Now().In(location).Format("2006-01-02"), "timezone": location.String(), "results": matched, "coverage": "Popular TMDB people sampled from the first 3 pages."})
}

func (h *CatalogController) fetch(c *gin.Context, path string, q url.Values, out any) bool {
	err := h.client.Get(c.Request.Context(), path, q, out)
	if err == nil {
		return true
	}
	var apiErr *tmdb.APIError
	if errors.As(err, &apiErr) {
		status := http.StatusBadGateway
		code := "tmdb_upstream_error"
		if apiErr.Status == http.StatusNotFound {
			status = http.StatusNotFound
			code = "catalog_item_not_found"
		}
		if apiErr.Status == http.StatusTooManyRequests {
			status = http.StatusServiceUnavailable
			code = "tmdb_rate_limited"
		}
		apiresponse.Error(c, status, code, "The catalogue provider could not complete this request")
		return false
	}
	apiresponse.Error(c, http.StatusBadGateway, "tmdb_unavailable", "The catalogue provider is temporarily unavailable")
	return false
}
func (h *CatalogController) respond(c *gin.Context, data any) {
	c.JSON(http.StatusOK, gin.H{"data": data, "meta": gin.H{"source": "TMDB", "image_base_url": h.imageBaseURL, "attribution": "This product uses the TMDB API but is not endorsed or certified by TMDB."}})
}
func listQuery(c *gin.Context) url.Values {
	page := c.DefaultQuery("page", "1")
	if n, err := strconv.Atoi(page); err != nil || n < 1 || n > 500 {
		page = "1"
	}
	return url.Values{"page": {page}, "region": {c.DefaultQuery("region", "DE")}}
}
func validNumericID(c *gin.Context) bool {
	if _, err := strconv.Atoi(c.Param("id")); err != nil {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_catalog_id", "Catalogue IDs must be numeric")
		return false
	}
	return true
}
func validNumber(c *gin.Context, key string) bool {
	if n, err := strconv.Atoi(c.Param(key)); err != nil || n < 0 {
		apiresponse.Error(c, http.StatusBadRequest, "invalid_"+key, key+" must be a non-negative number")
		return false
	}
	return true
}
