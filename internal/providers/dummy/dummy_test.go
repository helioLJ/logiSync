package dummy

import (
	"context"
	"testing"
)

func TestDummyTrack(t *testing.T) {
	provider := New("dummy")
	result, err := provider.Track(context.TODO(), "TEST123")
	if err != nil {
		t.Fatalf("track: %v", err)
	}

	if result.Payload["tracking_code"] != "TEST123" {
		t.Fatalf("unexpected tracking_code: %v", result.Payload["tracking_code"])
	}

	if len(result.Artifacts) != 1 {
		t.Fatalf("expected 1 artifact, got %d", len(result.Artifacts))
	}
}
