package storydna

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"plotlyst/server/tmdb"
)

const schemaVersion = 2

type Service struct {
	TMDB       *tmdb.Client
	Collection *mongo.Collection
}

func (s *Service) Profile(ctx context.Context, movieID int, includeSpoilers bool) (Profile, error) {
	if movieID <= 0 {
		return Profile{}, errors.New("a positive movie id is required")
	}
	var stored Profile
	if s.Collection.FindOne(ctx, bson.M{"movie_id": movieID, "provenance.schema_version": schemaVersion}).Decode(&stored) == nil {
		return redact(stored, includeSpoilers), nil
	}
	var movie tmdb.MovieDetail
	if err := s.TMDB.Get(ctx, fmt.Sprintf("/movie/%d", movieID), nil, &movie); err != nil {
		return Profile{}, err
	}
	profile := Classify(movie)
	_, err := s.Collection.UpdateOne(ctx, bson.M{"movie_id": movieID}, bson.M{"$set": profile}, options.UpdateOne().SetUpsert(true))
	if err != nil {
		return Profile{}, err
	}
	return redact(profile, includeSpoilers), nil
}

func (s *Service) Discover(ctx context.Context, req DiscoveryRequest) (Profile, []Match, []string, error) {
	req.Weights = normalizeWeights(req.Weights)
	seedID := req.SeedMovieID
	query := strings.TrimSpace(req.Query)
	if seedID == 0 && query != "" {
		var search tmdb.Page[tmdb.MovieSummary]
		if strings.Contains(strings.ToLower(query), "like ") {
			if err := s.TMDB.Get(ctx, "/search/movie", url.Values{"query": {seedTitle(query)}}, &search); err != nil {
				return Profile{}, nil, nil, err
			}
		} else {
			filters := url.Values{"sort_by": {"vote_count.desc"}}
			if floor := afterYear(query); floor > 0 {
				filters.Set("primary_release_date.gte", fmt.Sprintf("%d-01-01", floor+1))
			}
			lower := strings.ToLower(query)
			genreIDs := []string{}
			if strings.Contains(lower, "thriller") {
				genreIDs = append(genreIDs, "53")
			}
			if strings.Contains(lower, "psychological") || strings.Contains(lower, "mystery") {
				genreIDs = append(genreIDs, "9648")
			}
			if len(genreIDs) > 0 {
				filters.Set("with_genres", strings.Join(genreIDs, ","))
			}
			if err := s.TMDB.Get(ctx, "/discover/movie", filters, &search); err != nil {
				return Profile{}, nil, nil, err
			}
		}
		if len(search.Results) > 0 {
			seedID = search.Results[0].ID
		}
	}
	if seedID == 0 {
		return Profile{}, nil, nil, errors.New("choose a seed movie or describe a movie to compare")
	}
	seed, err := s.Profile(ctx, seedID, req.IncludeSpoilers)
	if err != nil {
		return Profile{}, nil, nil, err
	}

	var recommendations tmdb.Page[tmdb.MovieSummary]
	if err := s.TMDB.Get(ctx, fmt.Sprintf("/movie/%d/recommendations", seedID), url.Values{"page": {"1"}}, &recommendations); err != nil {
		return Profile{}, nil, nil, err
	}
	yearFloor := afterYear(query)
	wantSimpler := strings.Contains(strings.ToLower(query), "simpler") || strings.Contains(strings.ToLower(query), "less complex")
	wantNonlinear := strings.Contains(strings.ToLower(query), "nonlinear") || strings.Contains(strings.ToLower(query), "non-linear")
	matches := make([]Match, 0, len(recommendations.Results))
	for _, candidate := range recommendations.Results {
		if candidate.ID == seedID {
			continue
		}
		profile := ClassifySummary(candidate)
		if yearFloor > 0 && profile.ReleaseYear <= yearFloor {
			continue
		}
		if wantSimpler && profile.Scores.Complexity >= seed.Scores.Complexity {
			continue
		}
		if wantNonlinear && !contains(profile.Structures, "potentially-nonlinear") {
			continue
		}
		match := Compare(seed, profile, req.Weights)
		match.Profile = redact(profile, req.IncludeSpoilers)
		matches = append(matches, match)
		_, _ = s.Collection.UpdateOne(ctx, bson.M{"movie_id": profile.MovieID}, bson.M{"$set": profile}, options.UpdateOne().SetUpsert(true))
	}
	sort.SliceStable(matches, func(i, j int) bool {
		if matches[i].Similarity == matches[j].Similarity {
			return matches[i].MovieID < matches[j].MovieID
		}
		return matches[i].Similarity > matches[j].Similarity
	})
	if len(matches) > 12 {
		matches = matches[:12]
	}
	limitations := []string{"Story DNA is a transparent metadata heuristic, not a statement of artistic fact.", "TMDB catalogue text is not sent to an AI provider.", "Ending characteristics stay hidden unless an authorized classification supplies them and spoilers are requested."}
	return seed, matches, limitations, nil
}

func Classify(movie tmdb.MovieDetail) Profile {
	ids := make([]int, 0, len(movie.Genres))
	for _, genre := range movie.Genres {
		ids = append(ids, genre.ID)
	}
	return classify(movie.ID, movie.Title, movie.PosterPath, year(movie.ReleaseDate), movie.Runtime, ids, movie.VoteAverage)
}
func ClassifySummary(movie tmdb.MovieSummary) Profile {
	return classify(movie.ID, movie.Title, movie.PosterPath, year(movie.ReleaseDate), 0, movie.GenreIDs, movie.VoteAverage)
}

func classify(id int, title, poster string, releaseYear, runtime int, genres []int, vote float64) Profile {
	has := func(id int) bool {
		for _, g := range genres {
			if g == id {
				return true
			}
		}
		return false
	}
	complexity, tension, darkness, warmth, hope, pacing := 38, 42, 28, 48, 55, 52
	structures := []string{"linear"}
	styles := []string{"character-led"}
	tones := []string{"balanced"}
	if has(9648) || has(53) {
		complexity += 22
		tension += 30
		styles = append(styles, "mystery-driven")
	}
	if has(878) || has(14) {
		complexity += 16
		styles = append(styles, "world-building")
	}
	if has(80) || has(27) {
		darkness += 32
		tension += 18
		warmth -= 18
		tones = []string{"dark", "tense"}
	}
	if has(35) || has(10749) {
		warmth += 30
		hope += 18
		tones = append(tones, "warm")
	}
	if has(28) || has(12) {
		pacing += 28
		tension += 10
		styles = append(styles, "kinetic")
	}
	if has(18) {
		pacing -= 12
		styles = append(styles, "dialogue-led")
	}
	if runtime > 145 {
		complexity += 9
		pacing -= 8
	}
	if has(9648) && (has(878) || has(53)) {
		structures = []string{"layered", "potentially-nonlinear"}
	}
	visual := 48
	if has(878) || has(14) || has(16) || has(12) {
		visual = 76
	}
	dialogue := 52
	if has(18) || has(80) {
		dialogue = 68
	}
	now := time.Now().UTC()
	return Profile{MovieID: id, Title: title, PosterPath: poster, ReleaseYear: releaseYear, Structures: unique(structures), Styles: unique(styles), Tones: unique(tones), Scores: Scores{clamp(warmth), clamp(tension), clamp(darkness), clamp(hope), clamp(pacing), clamp(complexity), clamp(visual), clamp(dialogue)}, Spoiler: SpoilerTraits{EndingShape: "unclassified", Resolution: "unclassified"}, Provenance: Provenance{Method: "deterministic-metadata-heuristic-v1", Source: "TMDB genre, runtime, release and vote metadata", Policy: "No TMDB plot text sent to AI; interpret as an inferred discovery aid.", Confidence: confidence(runtime, len(genres), vote), GeneratedAt: now, SchemaVersion: schemaVersion}, UpdatedAt: now}
}

func Compare(a, b Profile, w Weights) Match {
	mood := meanDistance([]int{a.Scores.EmotionalWarmth, a.Scores.Tension, a.Scores.Darkness, a.Scores.Hope}, []int{b.Scores.EmotionalWarmth, b.Scores.Tension, b.Scores.Darkness, b.Scores.Hope})
	pace := 1 - math.Abs(float64(a.Scores.Pacing-b.Scores.Pacing))/100
	complexity := 1 - math.Abs(float64(a.Scores.Complexity-b.Scores.Complexity))/100
	structure := jaccard(a.Structures, b.Structures)
	score := 100 * (mood*w.Mood + pace*w.Pacing + complexity*w.Complexity + structure*w.Structure)
	reasons := []string{}
	if mood >= .75 {
		reasons = append(reasons, "A closely matched emotional contour")
	}
	if pace >= .78 {
		reasons = append(reasons, "A similar narrative tempo")
	}
	if complexity >= .78 {
		reasons = append(reasons, "Comparable story complexity")
	}
	if structure > 0 {
		reasons = append(reasons, "Shared "+strings.Join(intersection(a.Structures, b.Structures), ", ")+" structure")
	}
	if len(reasons) == 0 {
		reasons = append(reasons, "The strongest available balance across your selected characteristics")
	}
	return Match{MovieID: b.MovieID, Title: b.Title, PosterPath: b.PosterPath, ReleaseYear: b.ReleaseYear, Similarity: math.Round(score*10) / 10, Reasons: reasons}
}

func normalizeWeights(w Weights) Weights {
	if w.Mood < 0 {
		w.Mood = 0
	}
	if w.Pacing < 0 {
		w.Pacing = 0
	}
	if w.Complexity < 0 {
		w.Complexity = 0
	}
	if w.Structure < 0 {
		w.Structure = 0
	}
	total := w.Mood + w.Pacing + w.Complexity + w.Structure
	if total == 0 {
		return Weights{.3, .2, .25, .25}
	}
	w.Mood /= total
	w.Pacing /= total
	w.Complexity /= total
	w.Structure /= total
	return w
}
func redact(p Profile, include bool) Profile {
	if !include {
		p.Spoiler = SpoilerTraits{}
	}
	return p
}
func clamp(v int) int {
	if v < 0 {
		return 0
	}
	if v > 100 {
		return 100
	}
	return v
}
func confidence(runtime, count int, vote float64) float64 {
	c := .42 + float64(min(count, 5))*.035
	if runtime > 0 {
		c += .08
	}
	if vote > 0 {
		c += .03
	}
	return math.Round(math.Min(c, .68)*100) / 100
}
func year(raw string) int {
	if len(raw) < 4 {
		return 0
	}
	n, _ := strconv.Atoi(raw[:4])
	return n
}
func unique(v []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, s := range v {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	sort.Strings(out)
	return out
}
func intersection(a, b []string) []string {
	set := map[string]bool{}
	for _, s := range a {
		set[s] = true
	}
	out := []string{}
	for _, s := range b {
		if set[s] {
			out = append(out, s)
		}
	}
	return out
}
func contains(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func jaccard(a, b []string) float64 {
	union := map[string]bool{}
	for _, s := range a {
		union[s] = true
	}
	for _, s := range b {
		union[s] = true
	}
	if len(union) == 0 {
		return 1
	}
	return float64(len(intersection(a, b))) / float64(len(union))
}
func meanDistance(a, b []int) float64 {
	total := 0.0
	for i := range a {
		total += 1 - math.Abs(float64(a[i]-b[i]))/100
	}
	return total / float64(len(a))
}

var afterPattern = regexp.MustCompile(`(?i)(?:after|since)\s+(19\d{2}|20\d{2})`)

func afterYear(q string) int {
	m := afterPattern.FindStringSubmatch(q)
	if len(m) < 2 {
		return 0
	}
	n, _ := strconv.Atoi(m[1])
	return n
}
func seedTitle(q string) string {
	low := strings.ToLower(q)
	if i := strings.Index(low, "like "); i >= 0 {
		v := q[i+5:]
		for _, stop := range []string{" but ", " with ", " after ", " released "} {
			if j := strings.Index(strings.ToLower(v), stop); j >= 0 {
				v = v[:j]
			}
		}
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return q
}
