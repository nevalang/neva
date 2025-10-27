package funcs

import (
	"os"
	"strings"
	"testing"
)

func TestParseDotenv(t *testing.T) {
	content := strings.Join([]string{
		"# comment",
		"export FOO=bar",
		"BAR=some value",
		"BAZ='quoted#value'",
		"QUX=\"escaped\\nline\"",
		"TRIM= value with spaces  # trailing",
		"EMPTY=",
		"BOM=should be ignored",
	}, "\n")

	reader := strings.NewReader("\ufeff" + content)

	values, err := parseDotenv(reader)
	if err != nil {
		t.Fatalf("parseDotenv returned error: %v", err)
	}

	expected := map[string]string{
		"FOO":   "bar",
		"BAR":   "some value",
		"BAZ":   "quoted#value",
		"QUX":   "escaped\nline",
		"TRIM":  "value with spaces",
		"EMPTY": "",
		"BOM":   "should be ignored",
	}

	if len(values) != len(expected) {
		t.Fatalf("expected %d entries, got %d (%v)", len(expected), len(values), values)
	}

	for key, want := range expected {
		got, ok := values[key]
		if !ok {
			t.Fatalf("missing key %q", key)
		}
		if got != want {
			t.Fatalf("key %q: expected %q, got %q", key, want, got)
		}
	}
}

func TestParseDotenvErrors(t *testing.T) {
	cases := []string{
		"NOVALUE",
		"=novalue",
		"BROKEN='missing",
		"BAD=\"missing",
	}

	for _, input := range cases {
		_, err := parseDotenv(strings.NewReader(input))
		if err == nil {
			t.Fatalf("expected error for %q", input)
		}
	}
}

func TestLoadDotenvFile(t *testing.T) {
	t.Setenv("KEEP", "existing")

	file, err := os.CreateTemp(t.TempDir(), "dotenv-*.env")
	if err != nil {
		t.Fatalf("CreateTemp: %v", err)
	}
	defer file.Close()

	data := "KEEP=fromfile\nNEW=value\n"
	if _, err := file.WriteString(data); err != nil {
		t.Fatalf("WriteString: %v", err)
	}

	values, err := loadDotenvFile(file.Name())
	if err != nil {
		t.Fatalf("loadDotenvFile: %v", err)
	}

	if got := os.Getenv("NEW"); got != "value" {
		t.Fatalf("expected NEW=value in env, got %q", got)
	}

	if got := os.Getenv("KEEP"); got != "existing" {
		t.Fatalf("expected existing KEEP env untouched, got %q", got)
	}

	if values["KEEP"] != "fromfile" || values["NEW"] != "value" {
		t.Fatalf("unexpected parsed values: %v", values)
	}
}
