//go:build windows

package queue

import (
	"fmt"
	"os"
	"syscall"
	"unsafe"
)

var (
	modkernel32    = syscall.NewLazyDLL("kernel32.dll")
	procLockFileEx = modkernel32.NewProc("LockFileEx")
	procUnlockFile = modkernel32.NewProc("UnlockFile")
)

const lockfileExclusiveLock = 0x00000002

func (q *Queue) lock() (unlock func(), err error) {
	f, err := os.OpenFile(q.lockPath, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("open lock file: %w", err)
	}

	// LockFileEx(handle, flags, reserved, nNumberOfBytesToLockLow, nNumberOfBytesToLockHigh, overlapped)
	ol := new(syscall.Overlapped)
	r, _, e := procLockFileEx.Call(
		f.Fd(),
		lockfileExclusiveLock,
		0,
		1, 0,
		uintptr(unsafe.Pointer(ol)),
	)
	if r == 0 {
		f.Close()
		return nil, fmt.Errorf("LockFileEx: %w", e)
	}

	return func() {
		procUnlockFile.Call(f.Fd(), 0, 0, 1, 0)
		f.Close()
	}, nil
}
