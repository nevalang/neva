package app

import (
	"github.com/emil14/neva/internal/compiler"
	"github.com/emil14/neva/internal/compiler/coder"
	"github.com/emil14/neva/internal/compiler/parser"
	"github.com/emil14/neva/internal/compiler/program"
	"github.com/emil14/neva/internal/compiler/translator"
	"github.com/emil14/neva/internal/compiler/validator"
	"github.com/emil14/neva/internal/runtime"
	"github.com/emil14/neva/internal/runtime/connector"
	"github.com/emil14/neva/internal/runtime/decoder"
	"github.com/emil14/neva/internal/runtime/operators"
)

type App struct {
	compiler.Compiler
	dcdr decoder.Decoder
	rntm runtime.Runtime
}

func (app App) Run(bb []byte) (runtime.IO, error) {
	prog, err := app.dcdr.Decode(bb)
	if err != nil {
		return runtime.IO{}, err
	}
	return app.rntm.Run(prog)
}

func handleOnSend(msg runtime.Msg, from runtime.PortAddr) {
	// fmt.Printf("%v -> %v\n", from, msg)
}

func handleOnReceive(msg runtime.Msg, from, to runtime.PortAddr) {
	// fmt.Printf("%v <- %v <- %v\n", to, msg, from)
}

func MustNew() App {
	ops := program.NewOperators()

	return App{
		Compiler: compiler.MustNew(
			parser.MustNewYAML(),
			validator.New(),
			translator.New(ops),
			coder.New(),
			ops,
		),
		dcdr: decoder.MustNewJSON(),
		rntm: runtime.New(
			connector.MustNew(
				operators.New(),
				handleOnSend,
				handleOnReceive,
			),
		),
	}
}
