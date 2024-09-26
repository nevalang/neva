package golang

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/stretchr/testify/assert"
)

func TestGetPortChansMap(t *testing.T) {
	tests := []struct {
		name     string
		prog     *ir.Program
		expected map[ir.PortAddr]string
	}{
		{
			name: "Empty program",
			prog: &ir.Program{
				Ports:       map[ir.PortAddr]struct{}{},
				Connections: map[ir.PortAddr]ir.PortAddr{},
			},
			expected: map[ir.PortAddr]string{},
		},
		{
			name: "Single unconnected port",
			prog: &ir.Program{
				Ports: map[ir.PortAddr]struct{}{
					{Path: "main", Port: "in"}: {},
				},
				Connections: map[ir.PortAddr]ir.PortAddr{},
			},
			expected: map[ir.PortAddr]string{
				{Path: "main", Port: "in"}: "main_in",
			},
		},
		{
			name: "Two connected ports",
			prog: &ir.Program{
				Ports: map[ir.PortAddr]struct{}{
					{Path: "main", Port: "out"}:  {},
					{Path: "logger", Port: "in"}: {},
				},
				Connections: map[ir.PortAddr]ir.PortAddr{
					{Path: "main", Port: "out"}: {Path: "logger", Port: "in"},
				},
			},
			expected: map[ir.PortAddr]string{
				{Path: "main", Port: "out"}:  "main_out_to_logger_in",
				{Path: "logger", Port: "in"}: "main_out_to_logger_in",
			},
		},
		{
			name: "Multiple ports with array",
			prog: &ir.Program{
				Ports: map[ir.PortAddr]struct{}{
					{Path: "main", Port: "out", IsArray: true, Idx: 0}:  {},
					{Path: "main", Port: "out", IsArray: true, Idx: 1}:  {},
					{Path: "logger", Port: "in", IsArray: true, Idx: 0}: {},
					{Path: "logger", Port: "in", IsArray: true, Idx: 1}: {},
				},
				Connections: map[ir.PortAddr]ir.PortAddr{
					{Path: "main", Port: "out", IsArray: true, Idx: 0}: {Path: "logger", Port: "in", IsArray: true, Idx: 0},
					{Path: "main", Port: "out", IsArray: true, Idx: 1}: {Path: "logger", Port: "in", IsArray: true, Idx: 1},
				},
			},
			expected: map[ir.PortAddr]string{
				{Path: "main", Port: "out", IsArray: true, Idx: 0}:  "main_out_0_to_logger_in_0",
				{Path: "main", Port: "out", IsArray: true, Idx: 1}:  "main_out_1_to_logger_in_1",
				{Path: "logger", Port: "in", IsArray: true, Idx: 0}: "main_out_0_to_logger_in_0",
				{Path: "logger", Port: "in", IsArray: true, Idx: 1}: "main_out_1_to_logger_in_1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getPortChansMap(tt.prog)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestChanVarNameFromPortAddr(t *testing.T) {
	tests := []struct {
		addr     ir.PortAddr
		expected string
	}{
		{
			// empty
			addr:     ir.PortAddr{},
			expected: "_",
		},
		{
			// only path
			addr:     ir.PortAddr{Path: "logger/in"},
			expected: "logger_in_",
		},
		{
			// only port
			addr:     ir.PortAddr{Port: "a"},
			expected: "_a",
		},
		{
			// only port
			addr:     ir.PortAddr{Port: "a"},
			expected: "_a",
		},
		{
			// path and port
			addr:     ir.PortAddr{Path: "logger/in", Port: "a"},
			expected: "logger_in_a",
		},
		{
			// path, port and idx
			addr:     ir.PortAddr{Path: "logger/in", Port: "a", IsArray: true, Idx: 0},
			expected: "logger_in_a_0",
		},
		// idx is not 0 but IsArray is false
		{
			addr:     ir.PortAddr{Path: "logger/in", Port: "a", IsArray: false, Idx: 42},
			expected: "logger_in_a",
		},
	}
	for _, tt := range tests {
		t.Run(tt.addr.String(), func(t *testing.T) {
			result := chanVarNameFromPortAddr(tt.addr)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func Test_handleSpecialChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "const_ref_sender",
			input:    "$greeting",
			expected: "_greeting",
		},
		{
			name:     "normal port addr",
			input:    "foo:bar",
			expected: "foo_bar",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleSpecialChars(tt.input)
			if result != tt.expected {
				t.Errorf("handleSpecialChars() = %v, want %v", result, tt.expected)
			}
		})
	}
}
