package decoder

import (
	"github.com/emil14/respect/internal/runtime/program"
)

type Program struct {
	Operators map[string][]string `json:"operators"`
	Nodes     map[string]Node     `json:"nodes"`
	Net       map[string][]string `json:"net"`
}

type Node struct {
	Operator string              `json:"operator,omitempty"`
	In       map[string]PortMeta `json:"in"`
	Out      map[string]PortMeta `json:"out"`
}

type PortMeta struct {
	Slots uint8 `json:"slots"`
	Buf   uint8 `json:"buf"`
}

type Connection struct {
	From PortAddr   `json:"from"`
	To   []PortAddr `json:"to"`
}

type PortAddr struct {
	Node string `json:"node"`
	Port string `json:"port"`
	Idx  uint8  `json:"idx"`
}

type Decoder struct {
	unmarshal func([]byte, interface{}) error
	caster    interface {
		Cast(Program) program.Program
	}
}

func (d Decoder) Decode(bb []byte) (program.Program, error) {
	prog := Program{}
	if err := d.unmarshal(bb, &prog); err != nil {
		return program.Program{}, err
	}
	return d.caster.Cast(prog), nil
}

// func NewJSON() (Decoder, error) {
// 	return Decoder{
// 		unmarshal: json.Unmarshal,
// 		caster:    NewCaster(),
// 	}, nil
// }

// func MustNewJSON() Decoder {
// 	d, err := NewJSON()
// 	if err != nil {
// 		panic(err)
// 	}
// 	return d
// }
