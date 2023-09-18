package testutil_test

import (
	"testing"

	"github.com/dagger/graphql/testutil"
)

func TestSubsetSlice_Simple(t *testing.T) {

	super := []any{
		"1", "2", "3",
	}
	sub := []any{
		"3",
	}
	if !testutil.ContainSubsetSlice(super, sub) {
		t.Fatalf("expected slice to be subset of super, got false")
	}
}
func TestSubsetSlice_Simple_Fail(t *testing.T) {

	super := []any{
		"1", "2", "3",
	}
	sub := []any{
		"4",
	}
	if testutil.ContainSubsetSlice(super, sub) {
		t.Fatalf("expected slice to not be subset of super, got true")
	}
}
func TestSubsetSlice_NestedSlice(t *testing.T) {

	super := []any{
		[]any{
			"1", "2", "3",
		},
		[]any{
			"4", "5", "6",
		},
		[]any{
			"7", "8", "9",
		},
	}
	sub := []any{
		[]any{
			"2",
		},
		[]any{
			"9",
		},
		[]any{
			"5",
		},
	}
	if !testutil.ContainSubsetSlice(super, sub) {
		t.Fatalf("expected slice to be subset of super, got false")
	}
}
func TestSubsetSlice_NestedSlice_DifferentLength(t *testing.T) {

	super := []any{
		[]any{
			"1", "2", "3",
		},
		[]any{
			"4", "5", "6",
		},
		[]any{
			"7", "8", "9",
		},
	}
	sub := []any{
		[]any{
			"3",
		},
		[]any{
			"6",
		},
	}
	if !testutil.ContainSubsetSlice(super, sub) {
		t.Fatalf("expected slice to be subset of super, got false")
	}
}
func TestSubsetSlice_NestedSlice_Fail(t *testing.T) {

	super := []any{
		[]any{
			"1", "2", "3",
		},
		[]any{
			"4", "5", "6",
		},
		[]any{
			"7", "8", "9",
		},
	}
	sub := []any{
		[]any{
			"3",
		},
		[]any{
			"3",
		},
		[]any{
			"9",
		},
	}
	if !testutil.ContainSubsetSlice(super, sub) {
		t.Fatalf("expected slice to be subset of super, got false")
	}
}

func TestSubset_Simple(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	sub := map[string]any{
		"c": "3",
	}
	if !testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to be subset of super, got false")
	}

}
func TestSubset_Simple_Fail(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
	}
	sub := map[string]any{
		"d": "3",
	}
	if testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to not be subset of super, got true")
	}

}
func TestSubset_NestedMap(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": map[string]any{
			"aa": "11",
			"bb": "22",
			"cc": "33",
		},
	}
	sub := map[string]any{
		"c": "3",
		"d": map[string]any{
			"cc": "33",
		},
	}
	if !testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to be subset of super, got false")
	}
}
func TestSubset_NestedMap_Fail(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": map[string]any{
			"aa": "11",
			"bb": "22",
			"cc": "33",
		},
	}
	sub := map[string]any{
		"c": "3",
		"d": map[string]any{
			"dd": "44",
		},
	}
	if testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to not be subset of super, got true")
	}
}
func TestSubset_NestedSlice(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": []any{
			"11", "22",
		},
	}
	sub := map[string]any{
		"c": "3",
		"d": []any{
			"11",
		},
	}
	if !testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to be subset of super, got false")
	}
}
func TestSubset_ComplexMixed(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": map[string]any{
			"aa": "11",
			"bb": "22",
			"cc": []any{
				"ttt", "rrr", "sss",
			},
		},
		"e": []any{
			"111", "222", "333",
		},
		"f": []any{
			[]any{
				"9999", "8888", "7777",
			},
			[]any{
				"6666", "5555", "4444",
			},
		},
	}
	sub := map[string]any{
		"c": "3",
		"d": map[string]any{
			"bb": "22",
			"cc": []any{
				"sss",
			},
		},
		"e": []any{
			"111",
		},
		"f": []any{
			[]any{
				"8888", "9999",
			},
			[]any{
				"4444",
			},
		},
	}
	if !testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to be subset of super, got false")
	}
}
func TestSubset_ComplexMixed_Fail(t *testing.T) {

	super := map[string]any{
		"a": "1",
		"b": "2",
		"c": "3",
		"d": map[string]any{
			"aa": "11",
			"bb": "22",
			"cc": []any{
				"ttt", "rrr", "sss",
			},
		},
		"e": []any{
			"111", "222", "333",
		},
		"f": []any{
			[]any{
				"9999", "8888", "7777",
			},
			[]any{
				"6666", "5555", "4444",
			},
		},
	}
	sub := map[string]any{
		"c": "3",
		"d": map[string]any{
			"bb": "22",
			"cc": []any{
				"doesnotexist",
			},
		},
		"e": []any{
			"111",
		},
		"f": []any{
			[]any{
				"4444",
			},
		},
	}
	if testutil.ContainSubset(super, sub) {
		t.Fatalf("expected map to not be subset of super, got true")
	}
}
