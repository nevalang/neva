package funcs

import (
	"testing"

	"github.com/nevalang/neva/internal/runtime"
)

func TestStreamTagString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
		tag  streamTag
	}{
		{name: "open", tag: streamTagOpen, want: "Open"},
		{name: "data", tag: streamTagData, want: "Data"},
		{name: "close", tag: streamTagClose, want: "Close"},
		{name: "unknown", tag: streamTag(99), want: "streamTag(99)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if got := tt.tag.String(); got != tt.want {
				t.Fatalf("String() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStreamMessages(t *testing.T) {
	t.Parallel()

	data := runtime.NewStringMsg("value")
	openMsg := newStreamOpenMsg()
	dataMsg := newStreamDataMsg(data)
	closeMsg := newStreamCloseMsg()

	if !isStreamOpen(openMsg) || isStreamData(openMsg) || isStreamClose(openMsg) {
		t.Fatalf("open predicates mismatch: %v", openMsg)
	}
	if isStreamOpen(dataMsg) || !isStreamData(dataMsg) || isStreamClose(dataMsg) {
		t.Fatalf("data predicates mismatch: %v", dataMsg)
	}
	if isStreamOpen(closeMsg) || isStreamData(closeMsg) || !isStreamClose(closeMsg) {
		t.Fatalf("close predicates mismatch: %v", closeMsg)
	}
	if got := streamDataValue(dataMsg); !got.Equal(data) {
		t.Fatalf("streamDataValue() = %v, want %v", got, data)
	}
}

func TestStreamDataValuePanicsForNonData(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("streamDataValue did not panic")
		}
	}()

	streamDataValue(newStreamCloseMsg())
}
