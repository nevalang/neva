package disk

import (
	"context"
	"fmt"
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

func MustNew() Repository {
	return Repository{}
}
