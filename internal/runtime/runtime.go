package runtime

import (
	"errors"
)

var (
	ErrModNotFound = errors.New("module not found in env")
)
