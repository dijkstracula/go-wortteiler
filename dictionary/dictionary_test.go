package dictionary

import (
	"bufio"
	"strings"
	"testing"
)

func TestSetFromScanner(t *testing.T) {
	input := `A
	B
	C
	D`

	scanner := bufio.NewScanner(strings.NewReader(input))
	set, err := setFromScanner(scanner)
	if err != nil {
		t.Errorf("Error building set: %v\n", err)
	}

	// lol are you kidding
	var keys []string
	for k := range set {
		keys = append(keys, k)
	}
	if len(keys) != 4 {
		t.Errorf("Expected 4 keys, got %d\n", len(keys))
	}
}

func TestFromFiles(t *testing.T) {
	wordFile := "../db/de_words.txt"
	prefFile := "../db/de_suffixes.txt"
	suffFile := "../db/de_prefixes.txt"

	_, err := FromFiles(wordFile, prefFile, suffFile)
	if err != nil {
		t.Errorf("error building set: %v\n", err)
	}
}

func TestValidFunc(t *testing.T) {
	words := map[string]interface{}{
		"wort": struct{}{},
	}
	prefs := map[string]interface{}{
		"aus": struct{}{},
	}
	suffs := map[string]interface{}{
		"en": struct{}{},
	}

	d := Dictionary{words, prefs, suffs}
	f := ValidFunc(&d)

	// TODO: fairly sure I don't need to test UTF-8 chars here?
	tests := []struct {
		in string
		ok bool
	}{
		{"wort", true},
		{"aus", true},
		{"en", true},
		{"nichts", false},
	}

	for _, test := range tests {
		if actual := f(test.in); actual != test.ok {
			t.Errorf("f(%+q) produced %v, expected %v", test.in, actual, test.ok)
		}
	}
}
