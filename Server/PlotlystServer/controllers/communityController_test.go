package controllers

import "testing"

func TestValidMedia(t *testing.T) {
	for _, tc := range []struct {
		kind string
		id   int
		want bool
	}{{"movie", 1, true}, {"tv", 42, true}, {"person", 1, false}, {"movie", 0, false}} {
		if got := validMedia(tc.kind, tc.id); got != tc.want {
			t.Fatalf("validMedia(%q,%d)=%v want %v", tc.kind, tc.id, got, tc.want)
		}
	}
}

func TestCleanTextRemovesControlsAndLimitsRunes(t *testing.T) {
	if got := cleanText("  hello\x00 world  ", 100); got != "hello world" {
		t.Fatalf("unexpected sanitized text %q", got)
	}
	if got := cleanText("🎬cinema", 2); got != "🎬c" {
		t.Fatalf("expected rune-safe truncation, got %q", got)
	}
}
