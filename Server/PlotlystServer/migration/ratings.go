package migration

import "strings"

func LegacyRankingToStars(value int, name string) (float64, bool) {
	if value == 999 || strings.EqualFold(strings.TrimSpace(name), "Not_Ranked") {
		return 0, false
	}
	if value < 1 || value > 5 {
		return 0, false
	}
	return float64(6 - value), true
}
