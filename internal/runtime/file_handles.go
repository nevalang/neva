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

// FileHandles owns runtime file resources exposed to Neva as opaque integer handles.
type FileHandles struct {
	files          map[int64]*os.File
	stdioHandleIDs map[int64]struct{}
	nextID         int64
	mu             sync.RWMutex
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
func (handles *FileHandles) Add(file *os.File) int64 {
	handles.mu.Lock()
	defer handles.mu.Unlock()

	handleID := handles.nextID
	handles.nextID++
	handles.files[handleID] = file

	return handleID
}

// Get returns the file registered for handleID.
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
