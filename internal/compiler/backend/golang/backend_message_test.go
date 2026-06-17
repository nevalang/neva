package golang

import (
	"testing"

	"github.com/nevalang/neva/internal/compiler/ir"
)

func TestScalarListCtor(t *testing.T) {
	t.Parallel()

	got, ok := scalarListCtor([]ir.Message{
		{Type: ir.MsgTypeInt, Int: 1},
		{Type: ir.MsgTypeInt, Int: 2},
	})
	if !ok {
		t.Fatal("expected typed list ctor")
	}
	if got != "runtime.NewListIntMsg([]int64{1, 2})" {
		t.Fatalf("unexpected ctor: %s", got)
	}
}

func TestScalarListCtorMixedFallsBack(t *testing.T) {
	t.Parallel()

	_, ok := scalarListCtor([]ir.Message{
		{Type: ir.MsgTypeInt, Int: 1},
		{Type: ir.MsgTypeString, String: "x"},
	})
	if ok {
		t.Fatal("expected generic fallback for mixed list")
	}
}

func TestScalarDictCtor(t *testing.T) {
	t.Parallel()

	got, ok := scalarDictCtor(map[string]ir.Message{
		"a": {Type: ir.MsgTypeString, String: "x"},
		"b": {Type: ir.MsgTypeString, String: "y"},
	})
	if !ok {
		t.Fatal("expected typed dict ctor")
	}
	if got != "runtime.NewDictStringMsg(map[string]string{\"a\": \"x\", \"b\": \"y\"})" &&
		got != "runtime.NewDictStringMsg(map[string]string{\"b\": \"y\", \"a\": \"x\"})" {
		t.Fatalf("unexpected ctor: %s", got)
	}
}

func TestScalarDictCtorNonScalarFallsBack(t *testing.T) {
	t.Parallel()

	_, ok := scalarDictCtor(map[string]ir.Message{
		"a": {Type: ir.MsgTypeStruct},
	})
	if ok {
		t.Fatal("expected generic fallback for non-scalar dict")
	}
}
