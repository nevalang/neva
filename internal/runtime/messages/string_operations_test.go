package messages

import "testing"

func TestStringAtUsesRuneIndicesAndNegativeIndexes(t *testing.T) {
	t.Parallel()

	value := NewStringMsg("AЖ中😀Z")
	item, found := StringAt(value, -2)
	if !found || item.Str() != "😀" {
		t.Fatalf("StringAt(-2) = (%q, %t), want (😀, true)", item.Str(), found)
	}

	_, found = StringAt(value, 5)
	if found {
		t.Fatal("StringAt(5) found = true, want false")
	}
}

func TestStringSliceUsesRuneIndicesAndNormalizedBounds(t *testing.T) {
	t.Parallel()

	if got, want := StringSlice(NewStringMsg("AЖ中😀Z"), 1, 4).Str(), "Ж中😀"; got != want {
		t.Fatalf("StringSlice rune result = %q, want %q", got, want)
	}
	if got := StringSlice(NewStringMsg("abcdef"), 5, 2).Str(); got != "" {
		t.Fatalf("StringSlice reversed result = %q, want empty", got)
	}
}

func TestStringSplitUsesTypedStorage(t *testing.T) {
	t.Parallel()

	result := StringSplit(NewStringMsg("one,two"), NewStringMsg(","))
	values, ok := ListAsStrings(result.List())
	if !ok || len(values) != 2 || values[0] != "one" || values[1] != "two" {
		t.Fatalf("StringSplit result = %v, want typed [one two]", result)
	}
}
