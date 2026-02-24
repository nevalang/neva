package funcs

import "testing"

func TestLookupEnvResultMsg(t *testing.T) {
	t.Parallel()

	msg := lookupEnvResultMsg("value", true)
	if got := msg.Get("value").Str(); got != "value" {
		t.Fatalf("expected value field to be %q, got %q", "value", got)
	}

	if got := msg.Get("exists").Bool(); !got {
		t.Fatalf("expected exists field to be true")
	}
}
