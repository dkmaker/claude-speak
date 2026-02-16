package queue

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Queue manages a file-based FIFO message queue with file locking.
type Queue struct {
	dir      string
	filePath string
	lockPath string
}

// New creates a Queue that stores messages in the given directory.
func New(dir string) *Queue {
	return &Queue{
		dir:      dir,
		filePath: filepath.Join(dir, "queue"),
		lockPath: filepath.Join(dir, "queue.lock"),
	}
}

// Dir returns the queue directory path.
func (q *Queue) Dir() string {
	return q.dir
}

// Enqueue appends a message to the queue file atomically.
func (q *Queue) Enqueue(message string) error {
	if err := os.MkdirAll(q.dir, 0755); err != nil {
		return fmt.Errorf("create queue dir: %w", err)
	}

	unlock, err := q.lock()
	if err != nil {
		return fmt.Errorf("lock queue: %w", err)
	}
	defer unlock()

	f, err := os.OpenFile(q.filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("open queue: %w", err)
	}
	defer f.Close()

	// Strip newlines from message since queue is newline-delimited
	message = strings.ReplaceAll(message, "\n", " ")

	if _, err := fmt.Fprintln(f, message); err != nil {
		return fmt.Errorf("write queue: %w", err)
	}
	return nil
}

// Dequeue reads and removes the first message from the queue.
// Returns ("", false, nil) if queue is empty.
func (q *Queue) Dequeue() (string, bool, error) {
	unlock, err := q.lock()
	if err != nil {
		return "", false, fmt.Errorf("lock queue: %w", err)
	}
	defer unlock()

	f, err := os.Open(q.filePath)
	if os.IsNotExist(err) {
		return "", false, nil
	}
	if err != nil {
		return "", false, fmt.Errorf("open queue: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return "", false, nil
	}
	first := scanner.Text()

	// Read remaining lines
	var remaining []string
	for scanner.Scan() {
		remaining = append(remaining, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return "", false, fmt.Errorf("read queue: %w", err)
	}
	f.Close()

	// Rewrite queue without first line
	if len(remaining) == 0 {
		os.Remove(q.filePath)
	} else {
		if err := os.WriteFile(q.filePath, []byte(strings.Join(remaining, "\n")+"\n"), 0644); err != nil {
			return "", false, fmt.Errorf("rewrite queue: %w", err)
		}
	}

	return first, true, nil
}
