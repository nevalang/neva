package messages

import "testing"

func TestStreamMessageProtocol(t *testing.T) {
	t.Parallel()

	data := NewIntMsg(42)
	open := NewStreamOpenMsg()
	item := NewStreamDataMsg(data)
	closeMsg := NewStreamCloseMsg()

	if !IsStreamOpen(open) || IsStreamData(open) || IsStreamClose(open) {
		t.Fatal("Open event was classified incorrectly")
	}
	if IsStreamOpen(item) || !IsStreamData(item) || IsStreamClose(item) {
		t.Fatal("Data event was classified incorrectly")
	}
	if IsStreamOpen(closeMsg) || IsStreamData(closeMsg) || !IsStreamClose(closeMsg) {
		t.Fatal("Close event was classified incorrectly")
	}
	if got := StreamDataValue(item); !Equal(got, data) {
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
