package header

import "github.com/emil14/neva/internal/compiler/src"

type StdLib struct {
	lib        map[src.EntityRef]src.Component
	runtimeLib map[src.EntityRef]src.InterfaceDef
}

func (std StdLib) Entity(ref src.EntityRef) (src.Entity, error) {
	component, ok := std.lib[ref]
	if ok {
		return component, nil
	}

	std.runtimeLib[ref]

	return src.Component{}, nil
}
