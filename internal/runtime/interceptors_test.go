package runtime

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestJSONLPortSlotAddr_OmitsNilIndex(t *testing.T) {
	evt := jsonlSentEvent{
		Version:   jsonlTraceEventVersion,
		EventKind: jsonlEventSent,
		Index:     1,
		PortSlotAddr: PortSlotAddr{
			PortAddr: PortAddr{
				Path: "node",
				Port: "res",
			},
		},
		Msg: NewStringMsg("x"),
	}

	encoded, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("marshal jsonl sent event: %v", err)
	}
	if strings.Contains(string(encoded), "\"Index\":null") {
		t.Fatalf("expected nil index to be omitted, got: %s", string(encoded))
	}
}

func TestJSONLPortSlotAddr_EmitsArrayIndex(t *testing.T) {
	idx := uint8(2)
	evt := jsonlRecvEvent{
		Version: jsonlTraceEventVersion,
		Event:   jsonlEventRecv,
		Index:   1,
		PortSlotAddr: PortSlotAddr{
			PortAddr: PortAddr{
				Path: "node",
				Port: "args",
			},
			Index: &idx,
		},
		Msg: NewStringMsg("x"),
	}

	encoded, err := json.Marshal(evt)
	if err != nil {
		t.Fatalf("marshal jsonl recv event: %v", err)
	}
	if !strings.Contains(string(encoded), "\"Index\":2") {
		t.Fatalf("expected index to be present for array ports, got: %s", string(encoded))
	}
}
