package worker

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/dkmaker/claude-speak/internal/queue"
)

const defaultIdleTimeout = 60 * time.Second

// SpeakFunc processes a text message (calls API + plays audio).
type SpeakFunc func(text string) error

// Worker is the background TTS daemon that reads from a queue and speaks.
type Worker struct {
	queue       *queue.Queue
	speak       SpeakFunc
	idleTimeout time.Duration
	stopCh      chan struct{}
}

// Option configures a Worker.
type Option func(*Worker)

func WithIdleTimeout(d time.Duration) Option {
	return func(w *Worker) { w.idleTimeout = d }
}

// New creates a Worker. speakFn can be nil for testing (messages are discarded).
func New(q *queue.Queue, speakFn SpeakFunc, opts ...Option) *Worker {
	w := &Worker{
		queue:       q,
		speak:       speakFn,
		idleTimeout: defaultIdleTimeout,
		stopCh:      make(chan struct{}),
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

// Run starts the worker loop. Blocks until idle timeout or Stop is called.
func (w *Worker) Run() {
	// Write PID file
	pidPath := filepath.Join(w.queue.Dir(), "worker.pid")
	os.WriteFile(pidPath, []byte(strconv.Itoa(os.Getpid())), 0644)
	defer os.Remove(pidPath)

	idleStart := time.Now()

	for {
		select {
		case <-w.stopCh:
			return
		default:
		}

		msg, ok, err := w.queue.Dequeue()
		if err != nil {
			log.Printf("queue error: %v", err)
			time.Sleep(500 * time.Millisecond)
			continue
		}

		if !ok {
			// Queue empty — check idle timeout
			if time.Since(idleStart) >= w.idleTimeout {
				return
			}
			time.Sleep(250 * time.Millisecond)
			continue
		}

		// Reset idle timer
		idleStart = time.Now()

		if w.speak != nil {
			if err := w.speak(msg); err != nil {
				log.Printf("speak error: %v", err)
			}
		}
	}
}

// Stop signals the worker to exit.
func (w *Worker) Stop() {
	select {
	case <-w.stopCh:
	default:
		close(w.stopCh)
	}
}

// IsRunning checks if a worker process is running by reading the PID file.
func IsRunning(dir string) (int, bool) {
	pidPath := filepath.Join(dir, "worker.pid")
	data, err := os.ReadFile(pidPath)
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return 0, false
	}
	// Check if process exists
	proc, err := os.FindProcess(pid)
	if err != nil {
		return 0, false
	}
	if !isProcessRunning(proc) {
		return 0, false
	}
	return pid, true
}

// StopByPID reads the PID file and kills the worker process.
func StopByPID(dir string) error {
	pid, running := IsRunning(dir)
	if !running {
		return nil
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("find process %d: %w", pid, err)
	}
	return proc.Kill()
}
