package parser_test

import (
	"testing"

	ts "github.com/emil14/neva/pkg/types"
	h "github.com/emil14/neva/pkg/types/helper"
	"github.com/emil14/neva/pkg/types/parser"
	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	tests := []struct {
		s       string
		want    ts.Expr
		wantErr error
	}{
		// insts
		{
			s:    "t",
			want: h.Inst("t"),
		},
		{
			s:    " t ",
			want: h.Inst("t"),
		},
		{
			s:       "t y",
			wantErr: parser.ErrRefWIthSpace,
		},
		{
			s:       "t y<u8>",
			wantErr: parser.ErrRefWIthSpace,
		},
		{
			s:       "t<>",
			wantErr: parser.ErrEmptyAngleBrackets,
		},
		{
			s:       "t<",
			wantErr: parser.ErrMissingAngleClose,
		},
		{
			s:       "t<y",
			wantErr: parser.ErrMissingAngleClose,
		},
		{
			s:       "t<y<u8>",
			wantErr: parser.ErrInstArg,
		},
		{
			s:    "t<u8>",
			want: h.Inst("t", h.Inst("u8")),
		},
		{
			s:    "t <u8>",
			want: h.Inst("t", h.Inst("u8")),
		},
		{
			s:    " t <u8>",
			want: h.Inst("t", h.Inst("u8")),
		},
		{
			s:       "t<y<>>",
			wantErr: parser.ErrInstArg,
		},
		{
			s:    "t<y<u8>>",
			want: h.Inst("t", h.Inst("y", h.Inst("u8"))),
		},
		{
			s:       "t<y<u8>, u<>>",
			wantErr: parser.ErrInstArg,
		},
		{
			s: "t<y<u8>, u<u32>>",
			want: h.Inst(
				"t",
				h.Inst("y", h.Inst("u8")),
				h.Inst("u", h.Inst("u32")),
			),
		},
		// recs
		{
			s:       "{",
			wantErr: parser.ErrMissingCurlyClose,
		},
		{
			s:       "{ x",
			wantErr: parser.ErrMissingCurlyClose,
		},
		{
			s:       "{ x u8",
			wantErr: parser.ErrMissingCurlyClose,
		},
		{
			s:       "{ x u8, ",
			wantErr: parser.ErrMissingCurlyClose,
		},
		{
			s:    "{}",
			want: h.Rec(nil),
		},
		{
			s: "{ x u8 }",
			want: h.Rec(map[string]ts.Expr{
				"x": h.Inst("u8"),
			}),
		},
		{
			s:       "{ x u8 z }",
			wantErr: parser.ErrRecField,
		},
		{
			s:       "{ x u8<> }",
			wantErr: parser.ErrRecField,
		},
		{
			s:       "{ x u8, y }",
			wantErr: parser.ErrInvalidCurlyEl,
		},
		{
			s:       "{ x u8, z y x }",
			wantErr: parser.ErrRecField,
		},
		{
			s: "{ x u8, y u32 }",
			want: h.Rec(map[string]ts.Expr{
				"x": h.Inst("u8"),
				"y": h.Inst("u32"),
			}),
		},
		{
			s: "{ x t<u8> }",
			want: h.Rec(map[string]ts.Expr{
				"x": h.Inst("t", h.Inst("u8")),
			}),
		},
		{
			s: "{ x t<{y u8}> }",
			want: h.Rec(map[string]ts.Expr{
				"x": h.Inst(
					"t",
					h.Rec(map[string]ts.Expr{
						"y": h.Inst("u8"),
					}),
				),
			}),
		},
		// arrs
		{
			s:       "[",
			wantErr: parser.ErrBraceExprLen,
		},
		{
			s:       "[]",
			wantErr: parser.ErrBraceExprLen,
		},
		{
			s:       "[256u8",
			wantErr: parser.ErrMissingCloseBrace,
		},
		{
			s:       "[256]",
			wantErr: parser.ErrArrType,
		},
		{
			s:       "[256]t y",
			wantErr: parser.ErrArrType,
		},
		{
			s:       "[]u8",
			wantErr: parser.ErrArrSize,
		},
		{
			s:       "[abc]u8",
			wantErr: parser.ErrArrSize,
		},
		{
			s:    "[256]u8",
			want: h.Arr(256, h.Inst("u8")),
		},
		{
			s: "[256]{x u8}",
			want: h.Arr(256, h.Rec(map[string]ts.Expr{
				"x": h.Inst("u8"),
			})),
		},
		// Enums
		{
			s:    "{ x }",
			want: h.Enum("x"),
		},
		{
			s:    "{ x, y, z }",
			want: h.Enum("x", "y", "z"),
		},
		{
			s:       "{ x, y, z, }",
			wantErr: parser.ErrTrailingComma,
		},
		// Union
		// |
		// t |
		// t | y
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.s, func(t *testing.T) {
			got, err := parser.Parse(tt.s)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
