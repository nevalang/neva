package golang

// import (
// 	"strings"
// 	"testing"

// 	"github.com/nevalang/neva/internal/compiler/ir"
// 	"github.com/stretchr/testify/assert"
// )

// func TestGetPortChannels(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		prog     *ir.Program
// 		expected string
// 	}{
// 		{
// 			name: "Empty program",
// 			prog: &ir.Program{
// 				Ports:       map[ir.PortAddr]struct{}{},
// 				Connections: map[ir.PortAddr]map[ir.PortAddr]struct{}{},
// 			},
// 			expected: "var (\n)",
// 		},
// 		{
// 			name: "Single port",
// 			prog: &ir.Program{
// 				Ports: map[ir.PortAddr]struct{}{
// 					{Path: "worker/in", Port: "a"}: {},
// 				},
// 				Connections: map[ir.PortAddr]map[ir.PortAddr]struct{}{},
// 			},
// 			expected: "var (\n\tworker_in_a = make(chan runtime.OrderedMsg)\n)",
// 		},
// 		{
// 			name: "Connected ports with multiple ports per side",
// 			prog: &ir.Program{
// 				Ports: map[ir.PortAddr]struct{}{
// 					{Path: "main/in", Port: "a"}:       {},
// 					{Path: "main/in", Port: "b"}:       {},
// 					{Path: "worker/in", Port: "a"}:     {},
// 					{Path: "worker/in", Port: "b"}:     {},
// 					{Path: "worker/out", Port: "c"}:    {},
// 					{Path: "logger/in", Port: "a"}:     {},
// 					{Path: "logger/out", Port: "b"}:    {},
// 					{Path: "main/out", Port: "result"}: {},
// 				},
// 				Connections: map[ir.PortAddr]map[ir.PortAddr]struct{}{
// 					{Path: "main/in", Port: "a"}: {
// 						{Path: "worker/in", Port: "a"}: {},
// 					},
// 					{Path: "main/in", Port: "b"}: {
// 						{Path: "worker/in", Port: "b"}: {},
// 					},
// 					{Path: "worker/out", Port: "c"}: {
// 						{Path: "logger/in", Port: "a"}: {},
// 					},
// 					{Path: "logger/out", Port: "b"}: {
// 						{Path: "main/out", Port: "result"}: {},
// 					},
// 				},
// 			},
// 			expected: "var (\n\tmain_in_a = make(chan runtime.OrderedMsg)\n\tmain_in_b = make(chan runtime.OrderedMsg)\n\tworker_out_c = make(chan runtime.OrderedMsg)\n\tlogger_out_b = make(chan runtime.OrderedMsg)\n)",
// 		},
// 		{
// 			name: "Ports with slots",
// 			prog: &ir.Program{
// 				Ports: map[ir.PortAddr]struct{}{
// 					{Path: "main/in", Port: "a", IsArray: true, Idx: 0}:    {},
// 					{Path: "main/in", Port: "a", IsArray: true, Idx: 1}:    {},
// 					{Path: "worker/in", Port: "b", IsArray: true, Idx: 0}:  {},
// 					{Path: "worker/in", Port: "b", IsArray: true, Idx: 1}:  {},
// 					{Path: "worker/out", Port: "b", IsArray: true, Idx: 2}: {},
// 				},
// 				Connections: map[ir.PortAddr]map[ir.PortAddr]struct{}{
// 					{Path: "main/in", Port: "a", IsArray: true, Idx: 0}: {
// 						{Path: "worker/in", Port: "b", IsArray: true, Idx: 0}: {},
// 					},
// 					{Path: "main/in", Port: "a", IsArray: true, Idx: 1}: {
// 						{Path: "worker/in", Port: "b", IsArray: true, Idx: 1}: {},
// 					},
// 				},
// 			},
// 			expected: "var (\n\tmain_in_a_0 = make(chan runtime.OrderedMsg)\n\tmain_in_a_1 = make(chan runtime.OrderedMsg)\n\tworker_out_b_2 = make(chan runtime.OrderedMsg)\n)",
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := getPortChannels(tt.prog)
// 			assert.NoError(t, err)
// 			assert.Equal(t, tt.expected, strings.TrimSpace(result))
// 		})
// 	}
// }

// func TestGetChannelName(t *testing.T) {
// 	channelMap := map[ir.PortAddr]string{
// 		{Path: "main/in", Port: "a"}:                          "main_in_a",
// 		{Path: "worker/out", Port: "b"}:                       "worker_out_b",
// 		{Path: "logger/in", Port: "a", IsArray: true, Idx: 0}: "logger_in_a_0",
// 	}

// 	tests := []struct {
// 		name        string
// 		addr        ir.PortAddr
// 		expected    string
// 		expectError bool
// 	}{
// 		{
// 			name:        "Existing channel",
// 			addr:        ir.PortAddr{Path: "main/in", Port: "a"},
// 			expected:    "main_in_a",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Existing channel with slot",
// 			addr:        ir.PortAddr{Path: "logger/in", Port: "a", IsArray: true, Idx: 0},
// 			expected:    "logger_in_a_0",
// 			expectError: false,
// 		},
// 		{
// 			name:        "Non-existing channel",
// 			addr:        ir.PortAddr{Path: "unknown/in", Port: "x"},
// 			expected:    "",
// 			expectError: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			result, err := getChannelName(tt.addr, channelMap)
// 			if tt.expectError {
// 				assert.Error(t, err)
// 			} else {
// 				assert.NoError(t, err)
// 				assert.Equal(t, tt.expected, result)
// 			}
// 		})
// 	}
// }
