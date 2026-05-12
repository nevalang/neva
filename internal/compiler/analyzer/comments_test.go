package analyzer

import (
	"testing"

	src "github.com/nevalang/neva/pkg/ast"
	"github.com/nevalang/neva/pkg/core"
	"github.com/stretchr/testify/require"
)

func TestValidateEntityComments(t *testing.T) {
	baseIO := src.IO{
		In: map[string]src.Port{
			"start": {},
			"data":  {},
		},
		Out: map[string]src.Port{
			"res": {},
			"err": {},
		},
	}

	baseComments := &src.Comments{
		Inports: map[string]string{
			"start": "trigger",
			"data":  "payload",
		},
		Outports: map[string]string{
			"res": "result",
			"err": "error",
		},
		Meta: core.Meta{
			Location: core.Location{
				Filename: "main.neva",
			},
		},
	}

	tests := []struct {
		name        string
		wantErrPart string
		entity      src.Entity
	}{
		{
			name: "valid complete comments for interface",
			entity: src.Entity{
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: baseIO,
				},
				Comments: baseComments,
			},
		},
		{
			name: "type entity with port tags is rejected before pair check",
			entity: src.Entity{
				Kind: src.TypeEntity,
				Comments: &src.Comments{
					Inports: map[string]string{
						"data": "payload",
					},
					Meta: core.Meta{},
				},
			},
			wantErrPart: "allowed only on interface or component entities",
		},
		{
			name: "inports without outports are rejected",
			entity: src.Entity{
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: baseIO,
				},
				Comments: &src.Comments{
					Inports: map[string]string{
						"start": "trigger",
					},
					Meta: core.Meta{},
				},
			},
			wantErrPart: "must include both @inport and @outport",
		},
		{
			name: "unknown inport is rejected",
			entity: src.Entity{
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: baseIO,
				},
				Comments: &src.Comments{
					Inports: map[string]string{
						"unknown": "x",
						"data":    "ok",
					},
					Outports: map[string]string{
						"res": "ok",
						"err": "ok",
					},
					Meta: core.Meta{},
				},
			},
			wantErrPart: "unknown inport",
		},
		{
			name: "missing inport coverage is rejected",
			entity: src.Entity{
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: baseIO,
				},
				Comments: &src.Comments{
					Inports: map[string]string{
						"start": "trigger",
					},
					Outports: map[string]string{
						"res": "ok",
						"err": "ok",
					},
					Meta: core.Meta{},
				},
			},
			wantErrPart: "must document all inports",
		},
		{
			name: "missing outport coverage is rejected",
			entity: src.Entity{
				Kind: src.InterfaceEntity,
				Interface: src.Interface{
					IO: baseIO,
				},
				Comments: &src.Comments{
					Inports: map[string]string{
						"start": "trigger",
						"data":  "payload",
					},
					Outports: map[string]string{
						"res": "ok",
					},
					Meta: core.Meta{},
				},
			},
			wantErrPart: "must document all outports",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateEntityComments(&tt.entity)
			if tt.wantErrPart == "" {
				require.Nil(t, err)
				return
			}
			require.NotNil(t, err)
			require.Contains(t, err.Message, tt.wantErrPart)
		})
	}
}
