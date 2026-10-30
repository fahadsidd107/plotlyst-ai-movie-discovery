package debate

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"testing"
)

func TestGroundedObservationsAndPolarization(t *testing.T) {
	rows := []ReviewRecord{{ID: bson.NewObjectID(), Stars: 5, Body: "The acting and cinematography were extraordinary."}, {ID: bson.NewObjectID(), Stars: 4, Body: "Excellent acting from the cast."}, {ID: bson.NewObjectID(), Stars: 1, Body: "The pacing dragged and the writing felt weak."}, {ID: bson.NewObjectID(), Stars: 2, Body: "Slow pacing and muddled writing."}, {ID: bson.NewObjectID(), Stars: 3, Body: "A balanced film."}}
	loved, disliked, _, err := (DeterministicAnalyzer{}).Analyze(context.Background(), rows)
	if err != nil || len(loved) == 0 || len(disliked) == 0 {
		t.Fatalf("expected grounded observations: %#v %#v %v", loved, disliked, err)
	}
	p := Polarize(Distribution(rows))
	if !p.Polarized {
		t.Fatalf("expected polarized distribution: %#v", p)
	}
	for _, item := range append(loved, disliked...) {
		if len(item.Evidence) == 0 || item.Evidence[0].ReviewID == "" {
			t.Fatal("every observation must link evidence")
		}
	}
}
func TestSpoilersNeverEnterThemeSummary(t *testing.T) {
	rows := []ReviewRecord{{ID: bson.NewObjectID(), Stars: 5, Body: "The ending reveals the killer and acting", Spoiler: true}}
	loved, _, _, _ := (DeterministicAnalyzer{}).Analyze(context.Background(), rows)
	if len(loved) != 0 {
		t.Fatal("spoiler-tagged prose must not enter summaries")
	}
}
