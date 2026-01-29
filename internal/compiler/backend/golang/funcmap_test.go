package golang

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ir"
	"github.com/stretchr/testify/assert"
)

func TestGetPortChansMap(t *testing.T) {
	tests := []struct {
		name             string
		connections      map[ir.PortAddr]ir.PortAddr
		expectedMap      map[ir.PortAddr]string
		expectedVarNames []string
	}{
		{
			name:             "empty_program",
			connections:      map[ir.PortAddr]ir.PortAddr{},
			expectedMap:      map[ir.PortAddr]string{},
			expectedVarNames: []string{},
		},
		{
			name: "two_connected_ports",
			connections: map[ir.PortAddr]ir.PortAddr{
				{Path: "main", Port: "out"}: {Path: "logger", Port: "in"},
			},
			expectedMap: map[ir.PortAddr]string{
				{Path: "main", Port: "out"}:  "main_out_to_logger_in",
				{Path: "logger", Port: "in"}: "main_out_to_logger_in",
			},
			expectedVarNames: []string{"main_out_to_logger_in"},
		},
		{
			name: "multiple_ports_with_array",
			connections: map[ir.PortAddr]ir.PortAddr{
				{Path: "main", Port: "out", IsArray: true, Idx: 0}: {Path: "logger", Port: "in", IsArray: true, Idx: 0},
				{Path: "main", Port: "out", IsArray: true, Idx: 1}: {Path: "logger", Port: "in", IsArray: true, Idx: 1},
			},
			expectedMap: map[ir.PortAddr]string{
				{Path: "main", Port: "out", IsArray: true, Idx: 0}:  "main_out_0_to_logger_in_0",
				{Path: "main", Port: "out", IsArray: true, Idx: 1}:  "main_out_1_to_logger_in_1",
				{Path: "logger", Port: "in", IsArray: true, Idx: 0}: "main_out_0_to_logger_in_0",
				{Path: "logger", Port: "in", IsArray: true, Idx: 1}: "main_out_1_to_logger_in_1",
			},
			expectedVarNames: []string{
				"main_out_0_to_logger_in_0",
				"main_out_1_to_logger_in_1",
			},
		},
		{
			name: "single_ports_to_array_slots",
			connections: map[ir.PortAddr]ir.PortAddr{
				{Path: "read_file", Port: "res"}: {Path: "fan_in", Port: "data", IsArray: true, Idx: 0},
				{Path: "read_file", Port: "err"}: {Path: "fan_in", Port: "data", IsArray: true, Idx: 1},
			},
			expectedMap: map[ir.PortAddr]string{
				{Path: "read_file", Port: "res"}:                      "read_file_res_to_fan_in_data_0",
				{Path: "read_file", Port: "err"}:                      "read_file_err_to_fan_in_data_1",
				{Path: "fan_in", Port: "data", IsArray: true, Idx: 0}: "read_file_res_to_fan_in_data_0",
				{Path: "fan_in", Port: "data", IsArray: true, Idx: 1}: "read_file_err_to_fan_in_data_1",
			},
			expectedVarNames: []string{
				"read_file_res_to_fan_in_data_0",
				"read_file_err_to_fan_in_data_1",
			},
		},
	}

	b := Backend{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, _ := b.buildPortChanMap(tt.connections)
			assert.Equal(t, tt.expectedMap, result)
		})
	}
}

func TestChanVarNameFromPortAddr(t *testing.T) {
	tests := []struct { //nolint:govet // fieldalignment
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

	b := Backend{}

	for _, tt := range tests {
		t.Run(tt.addr.String(), func(t *testing.T) {
			result := b.chanVarNameFromPortAddr(tt.addr)
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
