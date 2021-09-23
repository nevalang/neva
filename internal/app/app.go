package app

import (
	"encoding/json"
	"fmt"

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

type Repo interface {
	GetModule(id string, major, minor, patch uint64) ([]byte, error)
}

type App struct {
	cmplr compiler.Compiler
	dcdr  decoder.Decoder
	rntm  runtime.Runtime
	repo  Repo
}

type programDescriptor struct {
	Deps map[string]map[string]struct {
		Repo    string `json:"repo"`
		Version string `json:"v"`
	} `json:"deps"`
	Import map[string]string `json:"import"`
	Root   string            `json:"root"`
}

func (app App) Compile(bb []byte) ([]byte, error) {
	pd := programDescriptor{}

	if err := json.Unmarshal(bb, &pd); err != nil {
		return nil, err
	}

	scope := map[string][]byte{}
	for alias, i := range pd.Import {

		scope[alias] = app.repo.GetModule()
	}
	// for every import get module

	app.cmplr.Compile()

	return nil, nil
}

func (app App) Run(bb []byte) (runtime.IO, error) {
	prog, err := app.dcdr.Decode(bb)
	if err != nil {
		return runtime.IO{}, err
	}
	return app.rntm.Run(prog)
}

func onSend(msg runtime.Msg, from runtime.PortAddr, to []runtime.PortAddr) runtime.Msg {
	fmt.Printf("%v -> %v\n", from, msg)
	return msg
}

func onReceive(msg runtime.Msg, from, to runtime.PortAddr) {
	fmt.Printf("%v <- %v <- %v\n", to, msg, from)
}

func MustNew() App {
	ops := program.NewOperators()
	cmplr := compiler.MustNew(
		parser.MustNewYAML(),
		validator.New(),
		translator.New(ops),
		coder.New(),
		ops,
	)
	dcdr := decoder.MustNewJSON()
	app := App{
		cmplr: cmplr,
		dcdr:  dcdr,
	}
	cnctr := connector.MustNew(operators.New(), onSend, onReceive)
	rntm := runtime.New(cnctr)
	app.rntm = rntm

	return app
}
