package main

import (
	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/analyze"
	"github.com/emil14/neva/internal/compiler/synth"
	rsrc "github.com/emil14/neva/internal/runtime/src"
)

func MustCreateCompiler() compiler.Compiler[rsrc.Program] {
	return compiler.MustNew[rsrc.Program](
		analyze.Analyzer{},
		synth.Synthesizer{},
	)
}
