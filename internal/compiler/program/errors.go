package program

import (
	"errors"
)

var (
	ErrModNotFound    = errors.New("module not found")
	ErrUnknownModType = errors.New("module has unknown type")
	ErrPortsLen       = errors.New("different number of ports")
	ErrPortTypes      = errors.New("different port types")
	ErrPortInvalid    = errors.New("invalid port")
	ErrPortNotFound   = errors.New("port not found")
)
