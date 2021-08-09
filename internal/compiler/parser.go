package compiler

import "github.com/emil14/stream/internal/core"

type Parser interface {
	Parse([]byte) (core.Component, error)
}
