package runtime

import "testing"

func TestStreamTagString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want string
		tag  StreamTag
	}{
		{name: "open", tag: StreamTagOpen, want: "Open"},
		{name: "data", tag: StreamTagData, want: "Data"},
		{name: "close", tag: StreamTagClose, want: "Close"},
		{name: "unknown", tag: StreamTag(99), want: "StreamTag(99)"},
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

	data := NewStringMsg("value")
	openMsg := NewStreamOpenMsg()
	dataMsg := NewStreamDataMsg(data)
	closeMsg := NewStreamCloseMsg()

	if !IsStreamOpen(openMsg) || IsStreamData(openMsg) || IsStreamClose(openMsg) {
		t.Fatalf("open predicates mismatch: %v", openMsg)
	}
	if IsStreamOpen(dataMsg) || !IsStreamData(dataMsg) || IsStreamClose(dataMsg) {
		t.Fatalf("data predicates mismatch: %v", dataMsg)
	}
	if IsStreamOpen(closeMsg) || IsStreamData(closeMsg) || !IsStreamClose(closeMsg) {
		t.Fatalf("close predicates mismatch: %v", closeMsg)
	}
	if got := StreamDataValue(dataMsg); !got.Equal(data) {
		t.Fatalf("StreamDataValue() = %v, want %v", got, data)
	}
}

func TestStreamDataValuePanicsForNonData(t *testing.T) {
	t.Parallel()

	defer func() {
		if recover() == nil {
			t.Fatal("StreamDataValue did not panic")
		}
	}()

	StreamDataValue(NewStreamCloseMsg())
}
