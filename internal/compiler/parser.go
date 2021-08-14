package compiler

import (
	"github.com/emil14/neva/internal/compiler/program"
)

type Parser interface {
	Parse([]byte) (program.Module, error)
	Unparse(program.Module) ([]byte, error)
}
