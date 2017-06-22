package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode"

	"github.com/dijkstracula/go-wortteiler/dictionary"
	"github.com/dijkstracula/go-wortteiler/splitter"
	"github.com/gorilla/mux"
)

var (
	logPrefix  = "[server]"
	reqTimeout = 5 * time.Second
)

func validateInput(s string) bool {
	if len(s) == 0 || len(s) > 64 {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

func splitFunc(s splitter.SplitFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		respCode := http.StatusOK
		var errString string
		var tree *splitter.Node

		ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
		defer cancel()

		logPrefix := fmt.Sprintf("%s %v ", logPrefix, r.RemoteAddr)

		// Write out valid json even on early returns
		defer func() {
			var blob []byte
			var err error

			// If we hit a snag, write out a json error blob.
			if respCode != http.StatusOK {
				errResp := make(map[string]string)
				errResp["err"] = errString

				blob, err = json.Marshal(errResp)
			} else {
				blob, err = json.Marshal(tree)
			}

			if err != nil {
				respCode = http.StatusInternalServerError
			}

			log.Printf("%s <- %d (%d bytes)\n", logPrefix, respCode, len(blob))
			w.WriteHeader(respCode)
			w.Write(blob)
		}()

		log.Printf("%s -> %s\n", logPrefix, r.URL.Path)

		// Grab the word, canonicalize, split, and translate.
		word, _ := mux.Vars(r)["word"]
		if !validateInput(word) {
			respCode = http.StatusBadRequest
			errString = "Invalid input word"
			return
		}

		word = strings.ToLower(word)
		tree = s(word)
		if err := dictionary.Translate(ctx, tree); err != nil {
			respCode = http.StatusInternalServerError
			errString = fmt.Sprintf("Translation error: %v", err)
			return
		}
	}
}

func notFound(w http.ResponseWriter, r *http.Request) {
	logPrefix := fmt.Sprintf("%s %v ", logPrefix, r.RemoteAddr)
	log.Printf("%s -> %s\n", logPrefix, r.URL.Path)
	log.Printf("%s <- %d\n", logPrefix, 404)
}

// New produces a new HTTP handler with the appropriate endpoints configured.
func New(splitter splitter.SplitFunc) http.Handler {
	r := mux.NewRouter()
	r.NotFoundHandler = http.HandlerFunc(notFound)

	r.NewRoute().
		Methods("GET").
		Path("/split/{word}").
		HandlerFunc(splitFunc(splitter))
	r.NewRoute().
		Methods("GET").
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("public/")))

	return r
}
