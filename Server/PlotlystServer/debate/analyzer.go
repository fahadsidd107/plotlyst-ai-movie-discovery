package debate

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"
)

type Analyzer interface {
	Analyze(context.Context, []ReviewRecord) ([]Observation, []Observation, string, error)
}

type DeterministicAnalyzer struct{}

var themes = []struct {
	Name  string
	Words []string
}{
	{"acting", []string{"acting", "actor", "actress", "performance", "performances", "cast"}},
	{"cinematography", []string{"cinematography", "camera", "visual", "visuals", "shot", "lighting"}},
	{"pacing", []string{"pacing", "pace", "slow", "rushed", "dragged", "tempo"}},
	{"writing", []string{"writing", "script", "screenplay", "dialogue", "story", "plot"}},
	{"soundtrack", []string{"soundtrack", "score", "music", "sound", "song"}},
	{"direction", []string{"direction", "director", "directing", "filmmaking"}},
}

func (DeterministicAnalyzer) Analyze(_ context.Context, rows []ReviewRecord) ([]Observation, []Observation, string, error) {
	return observations(rows, 4, 5, true), observations(rows, 1, 2, false), "deterministic-grounded-v1", nil
}

func observations(rows []ReviewRecord, minStars, maxStars int, positive bool) []Observation {
	type group struct {
		name     string
		evidence []Evidence
	}
	groups := map[string]*group{}
	for _, row := range rows {
		if row.Stars < minStars || row.Stars > maxStars || row.Spoiler {
			continue
		}
		lower := strings.ToLower(row.Body)
		for _, theme := range themes {
			matched := false
			for _, word := range theme.Words {
				if strings.Contains(lower, word) {
					matched = true
					break
				}
			}
			if !matched {
				continue
			}
			g := groups[theme.Name]
			if g == nil {
				g = &group{name: theme.Name}
				groups[theme.Name] = g
			}
			g.evidence = append(g.evidence, Evidence{ReviewID: row.ID.Hex(), AuthorName: fallbackName(row.AuthorName), Stars: row.Stars, Excerpt: excerpt(row.Body, 180), CreatedAt: row.CreatedAt})
		}
	}
	out := []Observation{}
	for _, g := range groups {
		if len(g.evidence) == 0 {
			continue
		}
		sentiment := "praised"
		if !positive {
			sentiment = "criticized"
		}
		out = append(out, Observation{Theme: g.name, Summary: fmt.Sprintf("%s was repeatedly %s by %d %s-rating reviewer%s.", title(g.name), sentiment, len(g.evidence), map[bool]string{true: "high", false: "low"}[positive], plural(len(g.evidence))), MentionCount: len(g.evidence), Evidence: g.evidence})
	}
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].MentionCount == out[j].MentionCount {
			return out[i].Theme < out[j].Theme
		}
		return out[i].MentionCount > out[j].MentionCount
	})
	if len(out) > 4 {
		out = out[:4]
	}
	return out
}

func Distribution(rows []ReviewRecord) map[string]int {
	out := map[string]int{"1": 0, "2": 0, "3": 0, "4": 0, "5": 0}
	for _, row := range rows {
		if row.Stars >= 1 && row.Stars <= 5 {
			out[fmt.Sprint(row.Stars)]++
		}
	}
	return out
}
func Polarize(distribution map[string]int) Polarization {
	total := 0
	for _, count := range distribution {
		total += count
	}
	if total == 0 {
		return Polarization{Summary: "Not enough ratings to measure disagreement."}
	}
	high := distribution["4"] + distribution["5"]
	low := distribution["1"] + distribution["2"]
	score := math.Round((2*float64(min(high, low))/float64(total))*100) / 100
	polarized := score >= .35 && high >= 2 && low >= 2
	summary := "Most ratings cluster on one side of the ballot."
	if polarized {
		summary = "High and low ratings both form meaningful camps; this measures disagreement, not quality."
	}
	return Polarization{Score: score, Polarized: polarized, Summary: summary}
}
func Filter(summary Summary, bucket string) Summary {
	allowed := func(stars int) bool {
		switch bucket {
		case "high":
			return stars >= 4
		case "low":
			return stars <= 2
		case "middle":
			return stars == 3
		default:
			return true
		}
	}
	filter := func(items []Observation) []Observation {
		out := []Observation{}
		for _, item := range items {
			evidence := []Evidence{}
			for _, entry := range item.Evidence {
				if allowed(entry.Stars) {
					evidence = append(evidence, entry)
				}
			}
			if len(evidence) > 0 {
				item.Evidence = evidence
				item.MentionCount = len(evidence)
				out = append(out, item)
			}
		}
		return out
	}
	summary.Loved = filter(summary.Loved)
	summary.Disliked = filter(summary.Disliked)
	return summary
}
func excerpt(raw string, max int) string {
	raw = strings.Join(strings.Fields(raw), " ")
	r := []rune(raw)
	if len(r) > max {
		return string(r[:max]) + "…"
	}
	return raw
}
func fallbackName(v string) string {
	if strings.TrimSpace(v) == "" {
		return "Plotlyst member"
	}
	return v
}
func title(v string) string {
	if v == "" {
		return v
	}
	return strings.ToUpper(v[:1]) + v[1:]
}
func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
