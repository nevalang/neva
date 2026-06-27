package funcs

import "testing"

func TestArgsListMsg(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		args []string
	}{
		{name: "empty", args: nil},
		{name: "one", args: []string{"neva"}},
		{name: "many", args: []string{"neva", "run", "main"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			list := argsListMsg(tt.args).List()
			if len(list) != len(tt.args) {
				t.Fatalf("expected %d args, got %d", len(tt.args), len(list))
			}

			for i, want := range tt.args {
				if got := list[i].Str(); got != want {
					t.Fatalf("arg[%d] = %q, want %q", i, got, want)
				}
			}
		})
	}
}
