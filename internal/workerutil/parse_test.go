package workerutil

import (
	"testing"

	"github.com/google/uuid"
)

func TestParseMessage(t *testing.T) {
	id := uuid.New()
	values := map[string]any{
		"jobId":        id.String(),
		"provider":     "dummy",
		"trackingCode": "TEST123",
	}

	jobID, provider, trackingCode, ok := ParseMessage(values)
	if !ok {
		t.Fatalf("expected ok")
	}
	if jobID != id {
		t.Fatalf("unexpected jobID")
	}
	if provider != "dummy" {
		t.Fatalf("unexpected provider")
	}
	if trackingCode != "TEST123" {
		t.Fatalf("unexpected tracking code")
	}
}

func TestParseMessageInvalid(t *testing.T) {
	values := map[string]any{
		"jobId":        "not-a-uuid",
		"provider":     "dummy",
		"trackingCode": "TEST123",
	}
	_, _, _, ok := ParseMessage(values)
	if ok {
		t.Fatalf("expected not ok")
	}
}
