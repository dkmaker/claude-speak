//go:build windows

package worker

import (
	"os"
	"syscall"
)

// isProcessRunning checks if a process with the given PID exists on Windows.
func isProcessRunning(proc *os.Process) bool {
	// On Windows, we can use GetExitCodeProcess via syscall
	// If the process is still running, it returns STILL_ACTIVE (259)
	handle, err := syscall.OpenProcess(syscall.PROCESS_QUERY_INFORMATION, false, uint32(proc.Pid))
	if err != nil {
		return false
	}
	defer syscall.CloseHandle(handle)

	var exitCode uint32
	err = syscall.GetExitCodeProcess(handle, &exitCode)
	if err != nil {
		return false
	}

	// STILL_ACTIVE = 259
	const STILL_ACTIVE = 259
	return exitCode == STILL_ACTIVE
}
