package compiler

import (
	cprog "github.com/emil14/stream/internal/compiler/program"
	rprog "github.com/emil14/stream/internal/runtime/program"
)

type Translator interface {
	Translate(cprog.Program) rprog.Program
}
