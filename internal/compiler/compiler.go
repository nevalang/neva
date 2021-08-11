package compiler

type compiler struct {
	parser Parser
}

func (c compiler) Compile(src []byte) ([]byte, error) {
	// srcMod, err := c.parser.Parse(src)
	// if err != nil {
	// 	return nil, err
	// }

	// if err := srcMod.Validate(); err != nil {
	// 	return nil, err
	// }

	// return program.Program{
	// 	Root: program.NodeMeta{},
	// }, nil

	return nil, nil
}

// func compileModule(mod core.Module) (program.Program, error) {

// }

// func New() Compiler
