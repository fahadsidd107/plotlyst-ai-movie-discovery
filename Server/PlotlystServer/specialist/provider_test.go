package specialist

import "testing"

func TestSourceStatusesRemainExplicit(t *testing.T) {
	for _, status := range []Status{StatusAvailable, StatusPartial, StatusUnavailable} {
		if status == "" {
			t.Fatal("status must not be empty")
		}
	}
}
func TestMediaReferenceShape(t *testing.T) {
	ref := MediaRef{MediaType: "movie", MediaID: 42}
	if ref.MediaType != "movie" || ref.MediaID != 42 {
		t.Fatal("media reference changed")
	}
}
