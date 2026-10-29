package storydna

import (
	"plotlyst/server/tmdb"
	"testing"
)

func TestClassificationIsBoundedAndTransparent(t *testing.T) {
	p := Classify(tmdb.MovieDetail{ID: 1, Title: "Example", Runtime: 150, ReleaseDate: "2024-01-01", Genres: []tmdb.Genre{{ID: 9648}, {ID: 878}}, VoteAverage: 8})
	if p.Scores.Complexity < 60 || p.Scores.Complexity > 100 {
		t.Fatalf("unexpected complexity %d", p.Scores.Complexity)
	}
	if p.Provenance.Method == "" || p.Provenance.Confidence >= 1 {
		t.Fatalf("heuristic provenance must be explicit: %#v", p.Provenance)
	}
}

func TestSimilarityWeightsChangeResult(t *testing.T) {
	a := Profile{Structures: []string{"linear"}, Scores: Scores{Pacing: 90, Complexity: 20}}
	b := Profile{MovieID: 2, Structures: []string{"layered"}, Scores: Scores{Pacing: 88, Complexity: 90}}
	pace := Compare(a, b, normalizeWeights(Weights{Pacing: 1}))
	complexity := Compare(a, b, normalizeWeights(Weights{Complexity: 1}))
	if pace.Similarity <= complexity.Similarity {
		t.Fatalf("pacing emphasis should favor this candidate: %.1f <= %.1f", pace.Similarity, complexity.Similarity)
	}
}

func TestQueryParsing(t *testing.T) {
	if got := seedTitle("Find movies like Inception but with a simpler storyline"); got != "Inception" {
		t.Fatalf("got %q", got)
	}
	if got := afterYear("nonlinear thrillers released after 2020"); got != 2020 {
		t.Fatalf("got %d", got)
	}
}
