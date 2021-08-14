package compiler

import (
	cprog "github.com/emil14/neva/internal/compiler/program"
	rprog "github.com/emil14/neva/internal/runtime/program"
)

type Translator interface {
	Translate(cprog.Program) rprog.Program
}
