package spoiler

import "testing"

func TestModesAndReviewConcealment(t *testing.T) {
	if Normalize("unknown") != ModeStandard {
		t.Fatal("unknown modes must fail safely")
	}
	if ConcealReview(ModeOff, true) {
		t.Fatal("off must reveal tagged reviews")
	}
	if !ConcealReview(ModeStandard, true) || ConcealReview(ModeStandard, false) {
		t.Fatal("standard must conceal tagged reviews only")
	}
	if !ConcealReview(ModeStrict, false) {
		t.Fatal("strict must conceal all community prose")
	}
}
