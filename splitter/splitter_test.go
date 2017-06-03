package splitter

import (
	"reflect"
	"testing"
)

func TestOffsets(t *testing.T) {
	tests := []struct {
		in  int
		out []int
	}{
		{
			0,
			nil,
		},
		{
			1,
			[]int{0},
		},
		{
			2,
			[]int{1, 0},
		},
		{
			3,
			[]int{1, 0, 2},
		},
		{
			4,
			[]int{2, 1, 3, 0},
		},
	}

	for _, test := range tests {
		if actual := offsets(test.in); !reflect.DeepEqual(actual, test.out) {
			t.Errorf("offsets(%d) expected %q, got %q", test.in, test.out, actual)
		}
	}
}

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
			[]partition{{"a", "bcde"}, {"abcd", "e"}, {"ab", "cde"}, {"abc", "de"}},
		},
	}

	for _, test := range tests {
		if actual := partitions(test.in); !reflect.DeepEqual(actual, test.out) {
			t.Errorf("partition(%q) expected %q, got %q", test.in, test.out, actual)
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
		expected *SplitNode
	}{
		{"x", nil},
		{"a", &SplitNode{"a", nil, nil}},
		{"ab",
			&SplitNode{"ab",
				&SplitNode{"a", nil, nil},
				&SplitNode{"b", nil, nil},
			},
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
		return str == "tag" ||
			str == "es" ||
			str == "buch" ||
			str == "tagesbuch"
	}

	splitter := Splitter(valid)
	tree := splitter("tagesbuch")

	if tree == nil {
		t.Errorf("Should have expected a valid tree; got nil")
	}
	if tree.word != "tagesbuch" {
		t.Errorf("Tree root should be \"tagesbuch\"; got %+q", tree.word)
	}

	// "tages / buch" is more even than "tag / esbuch", so make sure we get the former.
	if !reflect.DeepEqual(tree,
		&SplitNode{
			"tagesbuch",
			&SplitNode{"", &SplitNode{"tag", nil, nil}, &SplitNode{"es", nil, nil}},
			&SplitNode{"buch", nil, nil},
		}) {
		t.Errorf("Split not optmially even: got %+q", tree)
	}
}
