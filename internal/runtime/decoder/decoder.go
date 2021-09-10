package decoder

import (
	"encoding/json"

	"github.com/emil14/neva/internal/runtime/program"
)

type Program struct {
	RootNode NodeMeta             `json:"root"`
	Scope    map[string]Component `json:"scope"`
}

type Component struct {
	Operator string              `json:"operator,omitempty"`
	Workers  map[string]NodeMeta `json:"workers,omitempty"`
	Net      []Connection        `json:"net,omitempty"`
}

type NodeMeta struct {
	In        map[string]uint8 `json:"in"`
	Out       map[string]uint8 `json:"out"`
	Component string           `json:"component"`
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

func NewJSON() (Decoder, error) {
	return Decoder{
		unmarshal: json.Unmarshal,
		caster:    NewCaster(),
	}, nil
}

func MustNewJSON() Decoder {
	d, err := NewJSON()
	if err != nil {
		panic(err)
	}
	return d
}
