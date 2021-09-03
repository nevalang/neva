package validator

import (
	"testing"

	"github.com/emil14/neva/internal/compiler/program"
)

func Test_validator_validatePorts(t *testing.T) {
	t.Parallel()

	v := validator{}
	tests := []struct {
		name    string
		io      program.IO
		wantErr bool
	}{
		{
			name:    "empty ports",
			io:      program.IO{},
			wantErr: true,
		},
		{
			name: "invalid inports, no outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.UnknownType},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid inports, invalid outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.UnknownType},
				},
				Out: program.Ports{
					"x": {Type: program.UnknownType},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid inports, valid outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.UnknownType},
				},
				Out: program.Ports{
					"x": {Type: program.IntType},
				},
			},
			wantErr: true,
		},
		{
			name: "valid inports, no outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.IntType},
				},
			},
			wantErr: true,
		},
		{
			name: "valid inports, invalid outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.IntType},
				},
				Out: program.Ports{
					"x": {Type: program.UnknownType},
				},
			},
			wantErr: true,
		},
		{
			name: "valid inports, valid outports",
			io: program.IO{
				In: program.Ports{
					"x": {Type: program.IntType},
				},
				Out: program.Ports{
					"x": {Type: program.IntType},
				},
			},
			wantErr: false,
		},
	}

	for i := range tests {
		tt := tests[i]

		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if err := v.validatePorts(tt.io); (err != nil) != tt.wantErr {
				t.Errorf("validator.validatePorts() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
