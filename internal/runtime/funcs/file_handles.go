package funcs

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/nevalang/neva/internal/runtime"
)

//nolint:govet // Keeping map and counter together improves readability for this small state holder.
type fileHandleStore struct {
	mu     sync.Mutex
	nextID int64
	files  map[int64]*os.File
	stdio  map[int64]struct{}
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
		nextID: 3,
		files:  files,
		stdio: map[int64]struct{}{
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
	s.mu.Lock()
	defer s.mu.Unlock()

	file, found := s.files[handleID]
	if !found {
		return nil, fmt.Errorf("file handle %d not found", handleID)
	}

	return file, nil
}

func (s *fileHandleStore) Close(handleID int64) error {
	s.mu.Lock()
	if _, isStdio := s.stdio[handleID]; isStdio {
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
