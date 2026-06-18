package runtime

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

// Stdio file handle IDs mirror the conventional process file descriptors.
const (
	StdinFileHandleID  int64 = 0
	StdoutFileHandleID int64 = 1
	StderrFileHandleID int64 = 2
)

const fileHandleShardCount = 32

// FileHandles owns runtime file resources exposed to Neva as opaque integer handles.
type FileHandles struct {
	shards [fileHandleShardCount]fileHandleShard
	nextID atomic.Int64
}

type fileHandleShard struct {
	files map[int64]*os.File
	mu    sync.RWMutex
}

// NewFileHandles creates a runtime file-handle table with process stdio handles.
func NewFileHandles() *FileHandles {
	handles := &FileHandles{}
	handles.nextID.Store(StderrFileHandleID + 1)

	for shardID := range handles.shards {
		handles.shards[shardID].files = make(map[int64]*os.File)
	}

	handles.addStdio(StdinFileHandleID, os.Stdin)
	handles.addStdio(StdoutFileHandleID, os.Stdout)
	handles.addStdio(StderrFileHandleID, os.Stderr)

	return handles
}

// Add stores file and returns a new opaque runtime handle ID.
func (handles *FileHandles) Add(file *os.File) int64 {
	handleID := handles.nextID.Add(1) - 1
	shard := handles.shard(handleID)

	shard.mu.Lock()
	defer shard.mu.Unlock()

	shard.files[handleID] = file

	return handleID
}

// Get returns the file registered for handleID.
func (handles *FileHandles) Get(handleID int64) (*os.File, error) {
	shard := handles.shard(handleID)

	shard.mu.RLock()
	defer shard.mu.RUnlock()

	file, found := shard.files[handleID]
	if !found {
		return nil, fmt.Errorf("file handle %d not found", handleID)
	}

	return file, nil
}

// Close removes and closes a dynamic file handle.
func (handles *FileHandles) Close(handleID int64) error {
	if isStdioFileHandleID(handleID) {
		return fmt.Errorf("cannot close stdio file handle %d", handleID)
	}

	shard := handles.shard(handleID)

	shard.mu.Lock()
	file, found := shard.files[handleID]
	if found {
		delete(shard.files, handleID)
	}
	shard.mu.Unlock()

	if !found {
		return fmt.Errorf("file handle %d not found", handleID)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close file handle %d: %w", handleID, err)
	}
	return nil
}

func (handles *FileHandles) addStdio(handleID int64, file *os.File) {
	shard := handles.shard(handleID)
	shard.files[handleID] = file
}

func (handles *FileHandles) shard(handleID int64) *fileHandleShard {
	if handleID < 0 {
		return &handles.shards[0]
	}
	return &handles.shards[handleID%fileHandleShardCount]
}

func isStdioFileHandleID(handleID int64) bool {
	return handleID == StdinFileHandleID ||
		handleID == StdoutFileHandleID ||
		handleID == StderrFileHandleID
}
