package generator

import (
	"errors"
	"fbp/internal/parsing"
	"fbp/internal/runtime"
	"fbp/internal/types"
	"fmt"
)

type Translator struct {
	env map[string]runtime.Module
}

// Translate translates valid parsed module to runtime module.
func (t Translator) Translate(p parsing.Parsed) (runtime.Module, error) {
	deps := runtime.Deps{}
	for depname, depPorts := range p.Deps {
		v, ok := t.env[depname]
		if !ok {
			return runtime.Module{}, errors.New("unresolved dep")
		}

		if err := checkPorts(); err != nil {
			return runtime.Module{}, err
		}

		deps[depname] = v
	}

	inPorts := runtime.InPorts{}
	for portName, typeName := range p.In {
		t := types.ByName(typeName)
		if t == types.Unknown {
			return runtime.Module{}, fmt.Errorf("unknown type %s", typeName)
		}
		inPorts[portName] = t
	}

	outPorts := runtime.OutPorts{}
	for portName, typeName := range p.Out {
		t := types.ByName(typeName)
		if t == types.Unknown {
			return runtime.Module{}, fmt.Errorf("unknown type %s", typeName)
		}
		outPorts[portName] = t
	}

	return runtime.Module{
		Deps: deps,
		In:   inPorts,
		Out:  outPorts,
	}, nil
}

func checkPorts(parsing.Deps) error {
	return nil // TODO
}
