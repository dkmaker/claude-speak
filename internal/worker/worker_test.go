package worker_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/dkmaker/claude-speak/internal/queue"
	"github.com/dkmaker/claude-speak/internal/worker"
)

func TestWorkerWritesPID(t *testing.T) {
	dir := t.TempDir()
	q := queue.New(dir)

	w := worker.New(q, nil, worker.WithIdleTimeout(1*time.Second))
	go w.Run()
	defer w.Stop()

	time.Sleep(100 * time.Millisecond)

	pidFile := filepath.Join(dir, "worker.pid")
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		t.Error("worker.pid should exist while worker is running")
	}
}

func TestWorkerIdleTimeout(t *testing.T) {
	dir := t.TempDir()
	q := queue.New(dir)

	w := worker.New(q, nil, worker.WithIdleTimeout(500*time.Millisecond))

	done := make(chan struct{})
	go func() {
		w.Run()
		close(done)
	}()

	select {
	case <-done:
		// Worker exited due to idle timeout — good
	case <-time.After(3 * time.Second):
		t.Error("worker did not exit after idle timeout")
		w.Stop()
	}
}
