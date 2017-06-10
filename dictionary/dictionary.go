package dictionary

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"golang.org/x/net/dict"
)

var (
	dictdServer = flag.String("dictdServer", "all.dict.org:dict", "dictd server path")
	dictdDict   = flag.String("dictdDict", "fd-deu-eng", "dictd dictionary to lookup")

	// ErrWordNotFound is the error for when Translate() does not find a particular word.
	ErrWordNotFound = fmt.Errorf("No definitions found")
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
		t := strings.ToLower(scanner.Text())
		set[t] = struct{}{}
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

// Translate connects to the remote dictd instance and looks up the
// supplied word.  It's slightly annoying that this spins up and
// tears down a TCP connection for each lookup, but in prod we are
// running a localhost dictd instance and the responses are all
// small anyway, so a long-lived TCP connection is not going to be super
// helpful here.  Also, this means that we can round-robin around
// a load-balanced dictd server if we are not running on localhost.
//
// If the word is not found in the dictionary, returns sentinel error
// `ErrWordNotFound`.
func (d *Dictionary) Translate(deu string) (string, error) {
	client, err := dict.Dial("tcp", *dictdServer)
	if err != nil {
		return "", err
	}
	defer client.Close()

	defns, err := client.Define(*dictdDict, deu)
	if err != nil {
		return "", err
	}
	if len(defns) == 0 {
		return "", ErrWordNotFound
	}
	return defns[0].Word, nil
}

// ValidFunc produces a function that produces whether a given string is a
// valid word, prefix, or suffix.
func ValidFunc(d *Dictionary) func(string) bool {
	return func(s string) bool {
		// setFromScanner downcases all input strings so as to have
		// case-insensitive comparisons.
		s = strings.ToLower(s)

		// TODO: would be nice to short-circuit this, I guess
		_, isWord := d.Words[s]
		_, isPref := d.Prefixes[s]
		_, isSuff := d.Suffixes[s]
		return isWord || isPref || isSuff
	}
}
