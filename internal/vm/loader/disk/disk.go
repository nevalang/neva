package disk

import (
	"context"
	"fmt"
	"os"
)

type Loader struct{}

func (l Loader) Load(_ context.Context, pathToIRFile string) ([]byte, error) {
	irFile, err := os.ReadFile(pathToIRFile)
	if err != nil {
		return nil, fmt.Errorf("os read file: %w", err)
	}
	return irFile, nil
}

func MustNew() Loader {
	return Loader{}
}
