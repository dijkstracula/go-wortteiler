package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dijkstracula/go-wortteiler/dictionary"
	"github.com/dijkstracula/go-wortteiler/splitter"
	"github.com/gorilla/mux"
)

var (
	logPrefix  = "[server]"
	reqTimeout = 5 * time.Second
)

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
		word, ok := mux.Vars(r)["word"]
		if !ok {
			respCode = http.StatusBadRequest
			errString = "Missing 'word' parameter"
			return
		}
		if len(word) > 64 {
			respCode = http.StatusBadRequest
			errString = "Word too long"
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

// New produces a new HTTP handler with the appropriate endpoints configured.
func New(splitter splitter.SplitFunc) http.Handler {
	r := mux.NewRouter()

	r.NewRoute().
		Methods("POST").
		Path("/split/{word}").
		HandlerFunc(splitFunc(splitter))
	r.NewRoute().
		Methods("GET").
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("public/")))

	return r
}
