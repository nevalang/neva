package compiler

type (
	Parser     interface{ Parse([]byte) (Module, error) }
	Checker    interface{ Check(Program) error }
	Translator interface{ Translate(Program) (Module, error) }
)

type Compiler struct {
	parser     Parser
	checker    Checker
	translator Translator
	opsIO      map[OpRef]IO
}

func (c Compiler) PreCompile(Pkg) (Program, error) {
	return Program{}, nil
}

func (c Compiler) Compile(pkg Pkg) ([]byte, error) {
	if _, err := c.PreCompile(pkg); err != nil {
		return nil, err
	}

	return nil, nil
}
