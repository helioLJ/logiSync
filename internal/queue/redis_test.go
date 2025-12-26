package queue

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
)

func TestQueueLifecycle(t *testing.T) {
	mini, err := miniredis.Run()
	if err != nil {
		t.Skipf("miniredis unavailable: %v", err)
	}
	defer mini.Close()

	client := New(mini.Addr())
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := client.EnsureGroup(ctx, "tracking:jobs", "tracking-workers"); err != nil {
		t.Fatalf("ensure group: %v", err)
	}

	id, err := client.AddJob(ctx, "tracking:jobs", map[string]any{
		"jobId":        "11111111-1111-1111-1111-111111111111",
		"provider":     "dummy",
		"trackingCode": "TEST123",
	})
	if err != nil {
		t.Fatalf("add job: %v", err)
	}
	if id == "" {
		t.Fatalf("expected id")
	}

	messages, err := client.ReadGroup(ctx, "tracking:jobs", "tracking-workers", "worker-1", 1, 10*time.Millisecond)
	if err != nil {
		t.Fatalf("read group: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}

	if err := client.Ack(ctx, "tracking:jobs", "tracking-workers", messages[0].ID); err != nil {
		t.Fatalf("ack: %v", err)
	}
}
