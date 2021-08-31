package parser

import "testing"

func TestParser_parseCmd(t *testing.T) {
	tests := []struct {
		name string
		p    Parser
		raw  string
	}{
		{
			name: "",
			p:    Parser{},
			raw:  "set in.x int",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.parseCmd(tt.raw)
		})
	}
}
