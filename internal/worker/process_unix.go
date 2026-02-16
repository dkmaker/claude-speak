//go:build !windows

package worker

import (
	"os"
	"syscall"
)

// isProcessRunning checks if a process with the given PID exists on Unix systems.
func isProcessRunning(proc *os.Process) bool {
	// On Unix, FindProcess always succeeds. Send signal 0 to check if process exists.
	err := proc.Signal(syscall.Signal(0))
	return err == nil
}
