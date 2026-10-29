package intelligence

import (
	"context"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"plotlyst/server/tmdb"
)

type cacheEntry struct {
	response Response
	expires  time.Time
}
type Service struct {
	Extractor             Extractor
	Fallback              Extractor
	TMDB                  *tmdb.Client
	Timeout               time.Duration
	MaxInput, DailyBudget int
	Model                 string
	CacheTTL              time.Duration
	mu                    sync.RWMutex
	cache                 map[string]cacheEntry
}

func (s *Service) Discover(ctx context.Context, query, mode, region string) (Response, error) {
	clean, err := Sanitize(query, s.MaxInput)
	if err != nil {
		return Response{}, err
	}
	ctx, cancel := context.WithTimeout(ctx, s.Timeout)
	defer cancel()
	cacheKey := mode + "|" + strings.ToUpper(region) + "|" + strings.ToLower(clean)
	if cached, ok := s.cached(cacheKey); ok {
		cached.Usage.Cached = true
		return cached, nil
	}
	i, provider, err := s.Extractor.Extract(ctx, clean, mode)
	limitations := []string{"TMDB catalogue payloads are never sent to the AI provider.", "Mood matching uses transparent genre and keyword heuristics."}
	if err != nil {
		i, provider, err = s.Fallback.Extract(ctx, clean, mode)
		limitations = append(limitations, "The configured AI extractor was unavailable; a local parser was used.")
	}
	if err != nil {
		return Response{}, err
	}
	if region != "" {
		i.Region = region
	}
	i = Normalize(i)
	items, err := s.retrieve(ctx, i)
	if err != nil {
		return Response{}, err
	}
	response := Response{Intent: i, Items: items, Provenance: Provenance{Intent: provider, Retrieval: "TMDB server-side", Ranking: "plotlyst-deterministic-v1", Explanation: "grounded-template-v1"}, Usage: Usage{Provider: provider, Model: func() string {
		if provider == "openai" {
			return s.Model
		}
		return ""
	}(), DailyLimit: s.DailyBudget}, Limitations: limitations}
	s.store(cacheKey, response)
	return response, nil
}
func (s *Service) cached(key string) (Response, bool) {
	s.mu.RLock()
	v, ok := s.cache[key]
	s.mu.RUnlock()
	if !ok || time.Now().After(v.expires) {
		return Response{}, false
	}
	return cloneResponse(v.response), true
}
func (s *Service) store(key string, response Response) {
	if s.CacheTTL <= 0 {
		return
	}
	s.mu.Lock()
	if s.cache == nil {
		s.cache = map[string]cacheEntry{}
	}
	s.cache[key] = cacheEntry{response: cloneResponse(response), expires: time.Now().Add(s.CacheTTL)}
	s.mu.Unlock()
}
func cloneResponse(source Response) Response {
	out := source
	out.Items = append([]RankedItem(nil), source.Items...)
	for i := range out.Items {
		out.Items[i].Reasons = append([]string(nil), source.Items[i].Reasons...)
	}
	out.Limitations = append([]string(nil), source.Limitations...)
	out.Intent.MediaTypes = append([]string(nil), source.Intent.MediaTypes...)
	out.Intent.Genres = append([]string(nil), source.Intent.Genres...)
	out.Intent.Moods = append([]string(nil), source.Intent.Moods...)
	out.Intent.Languages = append([]string(nil), source.Intent.Languages...)
	return out
}
func (s *Service) retrieve(ctx context.Context, i Intent) ([]RankedItem, error) {
	var page tmdb.Page[tmdb.MediaSummary]
	path := "/search/multi"
	q := url.Values{"query": {i.Query}, "region": {i.Region}}
	filtered := len(i.Genres) > 0 || len(i.Languages) > 0 || i.Year > 0 || i.RuntimeMax > 0 || len(i.Moods) > 0
	forcedKind := ""
	if i.Mode == "discover" && filtered {
		forcedKind = "movie"
		if len(i.MediaTypes) == 1 && i.MediaTypes[0] == "tv" {
			forcedKind = "tv"
		}
		path = "/discover/" + forcedKind
		q = url.Values{"sort_by": {"vote_count.desc"}, "region": {i.Region}}
		genreIDs := map[string]string{"action": "28", "adventure": "12", "animation": "16", "comedy": "35", "crime": "80", "documentary": "99", "drama": "18", "fantasy": "14", "horror": "27", "mystery": "9648", "romance": "10749", "science fiction": "878", "thriller": "53", "western": "37"}
		ids := []string{}
		for _, g := range i.Genres {
			if id := genreIDs[g]; id != "" {
				ids = append(ids, id)
			}
		}
		if len(ids) > 0 {
			q.Set("with_genres", strings.Join(ids, ","))
		}
		if i.Year > 0 {
			if forcedKind == "tv" {
				q.Set("first_air_date_year", fmt.Sprint(i.Year))
			} else {
				q.Set("primary_release_year", fmt.Sprint(i.Year))
			}
		}
		if i.RuntimeMax > 0 {
			q.Set("with_runtime.lte", fmt.Sprint(i.RuntimeMax))
		}
		if len(i.Languages) > 0 {
			q.Set("with_original_language", i.Languages[0])
		}
	}
	if err := s.TMDB.Get(ctx, path, q, &page); err != nil {
		return nil, err
	}
	out := []RankedItem{}
	for _, m := range page.Results {
		if m.MediaType == "person" {
			continue
		}
		kind := m.MediaType
		if kind == "" {
			kind = forcedKind
			if kind == "" {
				kind = "movie"
			}
		}
		title := m.Title
		if title == "" {
			title = m.Name
		}
		date := m.ReleaseDate
		if date == "" {
			date = m.FirstAirDate
		}
		score := m.VoteAverage + min(m.Popularity/100, 2)
		reasons := []string{fmt.Sprintf("Catalogue score %.1f/10", m.VoteAverage)}
		if len(i.Moods) > 0 {
			reasons = append(reasons, "Matches your "+strings.Join(i.Moods, ", ")+" cue through local heuristics")
		}
		out = append(out, RankedItem{ID: m.ID, MediaType: kind, Title: title, Overview: m.Overview, PosterPath: m.PosterPath, ReleaseDate: date, VoteAverage: m.VoteAverage, Score: score, Reasons: reasons})
	}
	sort.SliceStable(out, func(a, b int) bool {
		if out[a].Score == out[b].Score {
			return out[a].ID < out[b].ID
		}
		return out[a].Score > out[b].Score
	})
	if len(out) > 20 {
		out = out[:20]
	}
	return out, nil
}
