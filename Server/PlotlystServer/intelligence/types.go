package intelligence

import (
	"errors"
	"regexp"
	"sort"
	"strings"
)

type Intent struct {
	Mode       string   `json:"mode"`
	Query      string   `json:"query"`
	MediaTypes []string `json:"media_types"`
	Genres     []string `json:"genres"`
	Moods      []string `json:"moods"`
	Languages  []string `json:"languages"`
	Year       int      `json:"year,omitempty"`
	RuntimeMax int      `json:"runtime_max,omitempty"`
	Region     string   `json:"region"`
}
type Usage struct {
	Provider   string `json:"provider"`
	Model      string `json:"model,omitempty"`
	Cached     bool   `json:"cached"`
	DailyLimit int    `json:"daily_limit"`
}
type Provenance struct {
	Intent      string `json:"intent"`
	Retrieval   string `json:"retrieval"`
	Ranking     string `json:"ranking"`
	Explanation string `json:"explanation"`
}
type RankedItem struct {
	ID          int      `json:"id"`
	MediaType   string   `json:"media_type"`
	Title       string   `json:"title"`
	Overview    string   `json:"overview"`
	PosterPath  string   `json:"poster_path"`
	ReleaseDate string   `json:"release_date,omitempty"`
	VoteAverage float64  `json:"vote_average"`
	Score       float64  `json:"score"`
	Reasons     []string `json:"reasons"`
}
type Response struct {
	Intent      Intent       `json:"intent"`
	Items       []RankedItem `json:"items"`
	Provenance  Provenance   `json:"provenance"`
	Usage       Usage        `json:"usage"`
	Limitations []string     `json:"limitations"`
}

var controls = regexp.MustCompile(`[\x00-\x08\x0B\x0C\x0E-\x1F\x7F]`)

func Sanitize(raw string, max int) (string, error) {
	raw = strings.TrimSpace(controls.ReplaceAllString(raw, ""))
	if raw == "" {
		return "", errors.New("query is required")
	}
	r := []rune(raw)
	if len(r) > max {
		return "", errors.New("query is too long")
	}
	return raw, nil
}
func Normalize(i Intent) Intent {
	i.Region = strings.ToUpper(strings.TrimSpace(i.Region))
	if len(i.Region) != 2 {
		i.Region = "DE"
	}
	i.MediaTypes = allowed(i.MediaTypes, map[string]bool{"movie": true, "tv": true})
	if len(i.MediaTypes) == 0 {
		i.MediaTypes = []string{"movie", "tv"}
	}
	i.Moods = allowed(i.Moods, map[string]bool{"feel-good": true, "tense": true, "dark": true, "romantic": true, "funny": true, "thoughtful": true, "family": true})
	i.Genres = allowed(i.Genres, map[string]bool{"action": true, "adventure": true, "animation": true, "comedy": true, "crime": true, "documentary": true, "drama": true, "fantasy": true, "horror": true, "mystery": true, "romance": true, "science fiction": true, "thriller": true, "western": true})
	i.Languages = allowed(i.Languages, map[string]bool{"en": true, "de": true, "hi": true, "es": true, "fr": true, "ko": true, "ja": true, "it": true})
	if i.Year < 1888 || i.Year > 2100 {
		i.Year = 0
	}
	if i.RuntimeMax < 20 || i.RuntimeMax > 600 {
		i.RuntimeMax = 0
	}
	return i
}
func allowed(values []string, set map[string]bool) []string {
	out := []string{}
	seen := map[string]bool{}
	for _, v := range values {
		v = strings.ToLower(strings.TrimSpace(v))
		if set[v] && !seen[v] {
			out = append(out, v)
			seen[v] = true
		}
	}
	sort.Strings(out)
	return out
}
