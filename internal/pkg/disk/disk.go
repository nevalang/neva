package disk

import (
	"context"
	"fmt"
	"io/fs"
	"os"
)

type Repository struct{}

func (r Repository) ByPath(ctx context.Context, path string) ([]byte, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	return file, nil
}

func (r Repository) Save(path string, bb []byte) error {
	err := os.WriteFile(path, bb, fs.ModePerm)
	if err != nil {
		return fmt.Errorf("write file: %w", err)
	}
	return nil
}

func MustNew() Repository {
	return Repository{}
}
