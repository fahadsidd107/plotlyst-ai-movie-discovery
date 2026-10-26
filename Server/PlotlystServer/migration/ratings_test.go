package migration

import "testing"

func TestLegacyRankingToStars(t *testing.T) {
	tests := []struct {
		value int
		name  string
		stars float64
		ok    bool
	}{{1, "Excellent", 5, true}, {5, "Terrible\n", 1, true}, {999, "Not_Ranked", 0, false}, {8, "Unknown", 0, false}}
	for _, tt := range tests {
		stars, ok := LegacyRankingToStars(tt.value, tt.name)
		if stars != tt.stars || ok != tt.ok {
			t.Fatalf("LegacyRankingToStars(%d,%q)=(%v,%v)", tt.value, tt.name, stars, ok)
		}
	}
}
