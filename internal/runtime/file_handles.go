package runtime

import (
	"fmt"
	"os"
	"sync"
)

// Stdio file handle IDs mirror the conventional process file descriptors.
const (
	StdinFileHandleID  int64 = 0
	StdoutFileHandleID int64 = 1
	StderrFileHandleID int64 = 2
)

// FileHandles owns process files that Neva code addresses through opaque IDs.
//
// The generated Go runtime keeps one FileHandles table in the stdlib registry
// and shares it between all file-related extern functions. The table only
// protects handle lookup and lifecycle bookkeeping; actual file I/O happens on
// the returned *os.File after Get releases the lock.
//
// This intentionally starts with one simple table instead of sharding. Get uses
// a read lock, so concurrent reads/writes through already-open handles do not
// block each other at the table level. Add and Close take the write lock because
// they mutate the map and the next dynamic ID. If open/close churn becomes a
// measured bottleneck, sharding can be added inside this type without changing
// the public runtime API.
type FileHandles struct {
	// files maps Neva-visible handle IDs to live process files.
	files map[int64]*os.File
	// stdioHandleIDs marks process stdio handles owned by the host process.
	// Neva may read/write these handles, but must not close them; the process and
	// Go runtime own their lifetime.
	stdioHandleIDs map[int64]struct{}
	// nextID is a monotonic cursor for handles created by Open/Create. It is not
	// round-robin and IDs are not reused, which keeps stale handle use visible as
	// an error after Close removes the old entry.
	nextID int64
	// mu protects files, stdioHandleIDs, and nextID. Get takes RLock for a short
	// map lookup; Add and Close take Lock because they modify the table.
	mu sync.RWMutex
}

// NewFileHandles creates a runtime file-handle table with process stdio handles.
func NewFileHandles() *FileHandles {
	files := map[int64]*os.File{
		StdinFileHandleID:  os.Stdin,
		StdoutFileHandleID: os.Stdout,
		StderrFileHandleID: os.Stderr,
	}

	return &FileHandles{
		files:  files,
		nextID: StderrFileHandleID + 1,
		stdioHandleIDs: map[int64]struct{}{
			StdinFileHandleID:  {},
			StdoutFileHandleID: {},
			StderrFileHandleID: {},
		},
	}
}

// Add stores file and returns a new opaque runtime handle ID.
//
// The caller transfers lifecycle ownership to FileHandles. User code must close
// the returned ID through the file_close extern when it is done with the file.
func (handles *FileHandles) Add(file *os.File) int64 {
	handles.mu.Lock()
	defer handles.mu.Unlock()

	handleID := handles.nextID
	handles.nextID++
	handles.files[handleID] = file

	return handleID
}

// Get returns the file registered for handleID.
//
// The returned *os.File remains owned by FileHandles; callers may perform I/O on
// it but must not close it directly. The lookup lock is released before return,
// so long reads or writes do not hold the table lock.
func (handles *FileHandles) Get(handleID int64) (*os.File, error) {
	handles.mu.RLock()
	defer handles.mu.RUnlock()

	file, found := handles.files[handleID]
	if !found {
		return nil, fmt.Errorf("file handle %d not found", handleID)
	}

	return file, nil
}

// Close removes and closes a dynamic file handle.
//
// Close is intentionally not idempotent. A second close of the same dynamic ID
// reports an unknown handle, matching Go's own "use after close is an error"
// posture and making double-close bugs visible to Neva code through the err
// outport. Stdio handles also return an error because their lifetime belongs to
// the host process, not to the Neva program.
func (handles *FileHandles) Close(handleID int64) error {
	handles.mu.Lock()
	if _, isStdio := handles.stdioHandleIDs[handleID]; isStdio {
		handles.mu.Unlock()
		return fmt.Errorf("cannot close stdio file handle %d", handleID)
	}

	file, found := handles.files[handleID]
	if found {
		delete(handles.files, handleID)
	}
	handles.mu.Unlock()

	if !found {
		return fmt.Errorf("file handle %d not found", handleID)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close file handle %d: %w", handleID, err)
	}
	return nil
}
