package splitter

import (
	"reflect"
	"testing"
)

func TestPartitions(t *testing.T) {
	tests := []struct {
		in  string
		out []partition
	}{
		{
			"",
			nil,
		},
		{
			"a",
			nil,
		},
		{
			"ab",
			[]partition{{"a", "b"}},
		},
		{
			"abc",
			[]partition{{"a", "bc"}, {"ab", "c"}},
		},
		{
			"abcde",
			[]partition{{"a", "bcde"}, {"ab", "cde"}, {"abc", "de"}, {"abcd", "e"}},
		},
	}

	for _, test := range tests {
		if actual := partitions(test.in); !reflect.DeepEqual(actual, test.out) {
			t.Errorf("partition(%q) expected %q, got %q", test.in, test.out, actual)
		}
	}
}

func TestScore(t *testing.T) {
	tests := []struct {
		in  *Node
		num int
		den int
	}{
		{
			nil,
			0, (2 << 0) - 1,
		},
		{
			MakeLeaf("a"),
			1, (2 << 1) - 1,
		},
		{
			MakeNode("",
				MakeLeaf("a"),
				MakeLeaf("b")),
			2, (2 << 2) - 1,
		},
		{
			MakeNode("ab",
				MakeLeaf("a"),
				MakeLeaf("b")),
			3, (2 << 2) - 1,
		},
	}

	for _, test := range tests {
		if actualNum, actualDen := test.in.Score(); actualNum != test.num || actualDen != test.den {
			t.Errorf("%+q.Score() expected (%d,%d), got (%d,%d)", test.in, test.num, test.den, actualNum, actualDen)
		}
	}
}

func TestTrivialDicts(t *testing.T) {
	valid := func(str string) bool {
		return str == "a" || str == "b" || str == "ab"
	}

	splitter := Splitter(valid)

	tests := []struct {
		in       string
		expected *Node
	}{
		{"x", nil},
		{"a", MakeLeaf("a")},
		{"ab",
			MakeNode("ab",
				MakeLeaf("a"),
				MakeLeaf("b"),
			),
		},
	}

	for _, test := range tests {
		actual := splitter(test.in)
		if !reflect.DeepEqual(actual, test.expected) {
			t.Errorf("split(%q) expected %+q, got %+q", test.in, test.expected, actual)
		}
	}
}

func TestTagesbuch(t *testing.T) {
	valid := func(str string) bool {
		return str == "tages" ||
			str == "buch" ||
			str == "tagesbuch"
	}

	splitter := Splitter(valid)
	tree := splitter("tagesbuch")

	if tree == nil {
		t.Errorf("Should have expected a valid tree; got nil")
	}
	if tree.Word != "tagesbuch" {
		t.Errorf("Tree root should be \"tagesbuch\"; got %+q", tree.Word)
	}

	// "tages / buch" is more even than "tag / esbuch", so make sure we get the latter.
	if !reflect.DeepEqual(tree,
		MakeNode(
			"tagesbuch",
			MakeLeaf("tages"),
			MakeLeaf("buch"),
		)) {
		t.Errorf("Split not optmially even: got %+q", tree)
	}
}

func TestEntschuldigung(t *testing.T) {
	valid := func(str string) bool {
		return str == "ent" ||
			str == "schuld" ||
			str == "schuldig" ||
			str == "ig" ||
			str == "ung" ||
			str == "entschuldigung"
	}

	splitter := Splitter(valid)
	tree := splitter("entschuldigung")

	if tree == nil {
		t.Errorf("Should have expected a valid tree; got nil")
	}
	if tree.Word != "entschuldigung" {
		t.Errorf("Tree root should be \"entschuldigung\"; got %+q", tree.Word)
	}

	// "tages / buch" is more even than "tag / esbuch", so make sure we get the latter.
	if !reflect.DeepEqual(tree,
		MakeNode(
			"entschuldigung",
			MakeLeaf("ent"),
			MakeNode(
				"", /* entschuldigung */
				MakeNode(
					"schuldig", /* schuldig */
					MakeLeaf("schuld"),
					MakeLeaf("ig"),
				),
				MakeLeaf("ung"),
			),
		),
	) {
		t.Errorf("Split not optmially even: got %v", tree)
	}
}
