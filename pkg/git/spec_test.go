package git

import "testing"

func TestParseRepoSpec(t *testing.T) {
	t.Parallel()

	tests := map[string]struct {
		input        string
		wantLocation string
		wantRevision string
		wantCloneURL string
		wantIsLocal  bool
	}{
		"module path with revision": {
			input:        "github.com/example/template@v1.2.3",
			wantLocation: "github.com/example/template",
			wantRevision: "v1.2.3",
			wantCloneURL: "https://github.com/example/template",
		},
		"ssh": {
			input:        "git@github.com:nevalang/neva-template",
			wantLocation: "git@github.com:nevalang/neva-template",
			wantCloneURL: "git@github.com:nevalang/neva-template",
		},
		"ssh with revision": {
			input:        "git@github.com:nevalang/neva-template@main",
			wantLocation: "git@github.com:nevalang/neva-template",
			wantRevision: "main",
			wantCloneURL: "git@github.com:nevalang/neva-template",
		},
		"local absolute": {
			input:        "/tmp/template",
			wantLocation: "/tmp/template",
			wantCloneURL: "/tmp/template",
			wantIsLocal:  true,
		},
		"local relative": {
			input:        "../testdata/template",
			wantLocation: "../testdata/template",
			wantCloneURL: "../testdata/template",
			wantIsLocal:  true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			spec, err := ParseRepoSpec(tc.input)
			if err != nil {
				t.Fatalf("ParseRepoSpec() error = %v", err)
			}

			if spec.Location != tc.wantLocation {
				t.Fatalf("Location = %q, want %q", spec.Location, tc.wantLocation)
			}

			if spec.Revision != tc.wantRevision {
				t.Fatalf("Revision = %q, want %q", spec.Revision, tc.wantRevision)
			}

			if spec.CloneURL() != tc.wantCloneURL {
				t.Fatalf("CloneURL() = %q, want %q", spec.CloneURL(), tc.wantCloneURL)
			}

			if spec.IsLocal() != tc.wantIsLocal {
				t.Fatalf("IsLocal() = %v, want %v", spec.IsLocal(), tc.wantIsLocal)
			}
		})
	}
}

func TestParseRepoSpecErrors(t *testing.T) {
	t.Parallel()

	for _, input := range []string{"", "@rev", "   "} {
		t.Run(input, func(t *testing.T) {
			t.Parallel()
			if _, err := ParseRepoSpec(input); err == nil {
				t.Fatalf("ParseRepoSpec(%q) expected error", input)
			}
		})
	}
}
