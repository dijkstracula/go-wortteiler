package dictionary

import (
	"bufio"
	"os"
)

// Dictionary holds lookup tables for valid words, as well as
// additional words that are special in that they may connect
// two compound words together.  TODO: we will merge those
// connective words at some point to shorten the tree.
type Dictionary struct {
	Words    map[string]interface{}
	Prefixes map[string]interface{}
	Suffixes map[string]interface{}
}

func setFromScanner(scanner *bufio.Scanner) (map[string]interface{}, error) {
	set := make(map[string]interface{})

	for scanner.Scan() {
		set[scanner.Text()] = struct{}{}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return set, nil
}

func setFromFile(path string) (map[string]interface{}, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return setFromScanner(bufio.NewScanner(file))
}

// FromFiles produces a Dictionary given three paths to dictionary words.
func FromFiles(wordPath, prefixPath, suffixPath string) (*Dictionary, error) {
	words, err := setFromFile(wordPath)
	if err != nil {
		return nil, err
	}
	prefixes, err := setFromFile(prefixPath)
	if err != nil {
		return nil, err
	}
	suffixes, err := setFromFile(suffixPath)
	if err != nil {
		return nil, err
	}

	return &Dictionary{words, prefixes, suffixes}, nil
}

// ValidFunc produces a function that produces whether a given string is a
// valid word, prefix, or suffix.
func ValidFunc(d *Dictionary) func(string) bool {
	return func(s string) bool {
		// TODO: would be nice to short-circuit this, I guess
		_, isWord := d.Words[s]
		_, isPref := d.Prefixes[s]
		_, isSuff := d.Suffixes[s]
		return isWord || isPref || isSuff
	}
}
