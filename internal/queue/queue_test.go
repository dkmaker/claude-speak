package queue_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dkmaker/claude-speak/internal/queue"
)

func TestEnqueueAndDequeue(t *testing.T) {
	dir := t.TempDir()
	q := queue.New(dir)

	// Enqueue two messages
	if err := q.Enqueue("hello"); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	if err := q.Enqueue("world"); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}

	// Dequeue should return FIFO order
	msg, ok, err := q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue: %v", err)
	}
	if !ok || msg != "hello" {
		t.Errorf("got %q, want %q", msg, "hello")
	}

	msg, ok, err = q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue: %v", err)
	}
	if !ok || msg != "world" {
		t.Errorf("got %q, want %q", msg, "world")
	}

	// Empty queue
	_, ok, err = q.Dequeue()
	if err != nil {
		t.Fatalf("Dequeue: %v", err)
	}
	if ok {
		t.Error("expected empty queue")
	}
}

func TestQueueFileCreation(t *testing.T) {
	dir := t.TempDir()
	q := queue.New(dir)
	if err := q.Enqueue("test"); err != nil {
		t.Fatalf("Enqueue: %v", err)
	}
	queueFile := filepath.Join(dir, "queue")
	if _, err := os.Stat(queueFile); os.IsNotExist(err) {
		t.Error("queue file should exist after enqueue")
	}
}
