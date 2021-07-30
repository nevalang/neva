package core

import (
	"errors"
	"fmt"
)

var (
	ErrModNotFound    = errors.New("module not found")
	ErrUnknownModType = errors.New("module has unknown type")
	ErrPortsLen       = errors.New("different number of ports")
	ErrPortTypes      = errors.New("different port types")
	ErrPortInvalid    = errors.New("invalid port")
)

func errModNotFound(name string) error {
	return fmt.Errorf("%w: %s", ErrModNotFound, name)
}

func errUnknownModType(name string, mod Module) error {
	return fmt.Errorf("%w: mod '%s', type %T", ErrUnknownModType, name, mod)
}

func errPortTypes(want, got PortInterface) error {
	return fmt.Errorf("%w: want %s, got %s", ErrPortTypes, want, got)
}

func errPortsLen(want, got int) error {
	return fmt.Errorf("%w: want %d, got %d", ErrPortsLen, want, got)
}

func errPortInvalid(name string, err error) error {
	return fmt.Errorf("%w: port '%s', err %s", ErrPortInvalid, name, err)
}
