package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestDictValueByKeyTypedMiss(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		dict runtime.Msg
	}{
		{"bool", runtime.NewDictBoolMsg(map[string]bool{"present": true})},
		{"int", runtime.NewDictIntMsg(map[string]int64{"present": 1})},
		{"float", runtime.NewDictFloatMsg(map[string]float64{"present": 1.5})},
		{"string", runtime.NewDictStringMsg(map[string]string{"present": "x"})},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, found := dictValueByKey(tt.dict.Dict(), "missing"); found {
				t.Fatal("missing key was found")
			}
		})
	}
}
