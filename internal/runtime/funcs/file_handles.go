package funcs

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

// fileHandleStore owns runtime file resources exposed to Neva as opaque integer
// handles.
type fileHandleStore struct {
	// files maps each live handle ID to the Go file it represents.
	files map[int64]*os.File
	// stdioHandleIDs marks predefined process stdio handles that must not close.
	stdioHandleIDs map[int64]struct{}
	// nextID is the next dynamic handle ID allocated by Add.
	nextID int64
	// mu protects files, stdioHandleIDs, and nextID.
	mu sync.RWMutex
}

const (
	stdinFileHandleID  int64 = 0
	stdoutFileHandleID int64 = 1
	stderrFileHandleID int64 = 2
)

func newFileHandleStore() *fileHandleStore {
	files := map[int64]*os.File{
		stdinFileHandleID:  os.Stdin,
		stdoutFileHandleID: os.Stdout,
		stderrFileHandleID: os.Stderr,
	}

	return &fileHandleStore{
		nextID: stderrFileHandleID + 1,
		files:  files,
		stdioHandleIDs: map[int64]struct{}{
			stdinFileHandleID:  {},
			stdoutFileHandleID: {},
			stderrFileHandleID: {},
		},
	}
}

func (s *fileHandleStore) Add(file *os.File) int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	handleID := s.nextID
	s.nextID++
	s.files[handleID] = file

	return handleID
}

func (s *fileHandleStore) Get(handleID int64) (*os.File, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	file, found := s.files[handleID]
	if !found {
		return nil, fmt.Errorf("file handle %d not found", handleID)
	}

	return file, nil
}

func (s *fileHandleStore) Close(handleID int64) error {
	s.mu.Lock()
	if _, isStdio := s.stdioHandleIDs[handleID]; isStdio {
		s.mu.Unlock()
		return fmt.Errorf("cannot close stdio file handle %d", handleID)
	}

	file, found := s.files[handleID]
	if found {
		delete(s.files, handleID)
	}
	s.mu.Unlock()

	if !found {
		return fmt.Errorf("file handle %d not found", handleID)
	}

	if err := file.Close(); err != nil {
		return fmt.Errorf("close file handle %d: %w", handleID, err)
	}
	return nil
}

func fileHandleID(msg runtime.Msg) (int64, error) {
	idMsg, isIntMsg := msg.(runtime.IntMsg)
	if !isIntMsg {
		return 0, errors.New("file handle must be int")
	}

	return idMsg.Int(), nil
}
