package funcs

import "testing"

func TestArgsListMsg(t *testing.T) {
	t.Parallel()

	got := argsListMsg([]string{"neva", "run", "main"})
	list := got.List()

	if len(list) != 3 {
		t.Fatalf("expected 3 args, got %d", len(list))
	}

	if list[0].Str() != "neva" || list[1].Str() != "run" || list[2].Str() != "main" {
		t.Fatalf("unexpected args payload: %v", list)
	}
}
