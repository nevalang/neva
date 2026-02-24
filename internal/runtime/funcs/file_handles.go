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
}

func newFileHandleStore() *fileHandleStore {
	return &fileHandleStore{
		nextID: 1,
		files:  make(map[int64]*os.File),
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
