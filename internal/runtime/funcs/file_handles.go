package funcs

import (
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

	id := s.nextID
	s.nextID++
	s.files[id] = file

	return id
}

func (s *fileHandleStore) Get(id int64) (*os.File, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	file, ok := s.files[id]
	if !ok {
		return nil, fmt.Errorf("file handle %d not found", id)
	}

	return file, nil
}

func (s *fileHandleStore) Close(id int64) error {
	s.mu.Lock()
	if _, isStdio := s.stdio[id]; isStdio {
		s.mu.Unlock()
		return fmt.Errorf("cannot close stdio file handle %d", id)
	}

	file, ok := s.files[id]
	if ok {
		delete(s.files, id)
	}
	s.mu.Unlock()

	if !ok {
		return fmt.Errorf("file handle %d not found", id)
	}

	return file.Close()
}

func fileHandleID(msg runtime.Msg) (int64, error) {
	idMsg, ok := msg.(runtime.IntMsg)
	if !ok {
		return 0, fmt.Errorf("file handle must be int")
	}

	return idMsg.Int(), nil
}
