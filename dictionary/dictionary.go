package dictionary

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/dijkstracula/go-wortteiler/splitter"

	"golang.org/x/net/dict"
)

var (
	logPrefix = "[dictionary]"

	dictdServer = flag.String("dictdServer", "all.dict.org:dict", "dictd server path")
	dictdDict   = flag.String("dictdDict", "fd-deu-eng,german-english", "dictd dictionaries to lookup (comma-separated)")

	// ErrWordNotFound is the error for when Translate() does not find a particular word.
	ErrWordNotFound = fmt.Errorf("No definitions found")

	// ErrContextCanceled is the error produced when a tree walk is interrupted with a
	// context timeout or cancelation.
	ErrContextCanceled = fmt.Errorf("Context canceled during iteration")
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

// Translate iterates through a Tree and translates all the words in
// turn.
func Translate(ctx context.Context, tree *splitter.Node) error {
	var interrupted bool
	var globalErr error

	f := func(n *splitter.Node) {
		// Make this a no-op if the context is canceled, or if a previous
		// call to d.translate produced an error.
		if globalErr != nil {
			return
		}
		select {
		case <-ctx.Done():
			log.Printf("%s Context canceled: %v\n", logPrefix, ctx.Err())
			interrupted = true
			return
		default:
			break
		}
		defns, localErr := translate(n.Word)
		if localErr != nil && localErr != ErrWordNotFound {
			globalErr = localErr
		} else if len(defns) > 0 {
			n.Defns = defns
		}
	}

	tree.ForEach(f)

	if interrupted {
		globalErr = ErrContextCanceled
	}

	return globalErr
}

// translate connects to the remote dictd instance and looks up the
// supplied word.  It's slightly annoying that this spins up and
// tears down a TCP connection for each lookup, but in prod we are
// running a localhost dictd instance and the responses are all
// small anyway, so a long-lived TCP connection is not going to be super
// helpful here.  Also, this means that we can round-robin around
// a load-balanced dictd server if we are not running on localhost.
//
// If the word is not found in the dictionary, returns sentinel error
// `ErrWordNotFound`.
func translate(deu string) ([]string, error) {
	strset := make(map[string]bool)

	if len(deu) == 0 {
		return nil, nil
	}

	// TODO: make a worker pool of these rather than just restarting
	// the connection every time
	client, err := dict.Dial("tcp", *dictdServer)
	if err != nil {
		log.Printf("%s %v\n", logPrefix, err)
		return nil, err
	}
	defer client.Close()

	for _, dict := range strings.Split(*dictdDict, ",") {
		log.Printf("%s looking up %s in %s\n", logPrefix, deu, dict)
		defns, err := client.Define(dict, deu)
		if len(defns) == 0 {
			continue
		}
		if err != nil {
			log.Printf("%s %v\n", logPrefix, err)
			return nil, err
		}

		for _, d := range defns {
			// freedict deu-eng has multiple defns split on newlines.
			for _, line := range strings.Split(string(d.Text), "\n") {
				// german-english tends to separate multiple defns on a
				// a line with a semicolon.
				for _, s := range strings.Split(line, ";") {
					s = strings.TrimSpace(s)
					if len(s) == 0 {
						continue
					}
					// german-english also seems to have the word itself
					// in the translated response?
					if strings.EqualFold(s, deu) {
						continue
					}

					// german-english hands us pure dictionary
					// defs which are not super readable
					if strings.Contains(s,"{") ||
					   strings.Contains(s,"<") ||
					   strings.Contains(s,"/") ||
					   strings.Contains(s,"(") {
						continue 
					}
					strset[s] = true
				}
			}
		}
	}

	if len(strset) == 0 {
		return nil, ErrWordNotFound
	}

	var strs []string
	for s, _ := range(strset) {
		strs = append(strs, s)
	}
	
	sort.Slice(strs, func (i,j int) bool {
		return len(strs[i]) < len(strs[j])
	})

	return strs, nil
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
