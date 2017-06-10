package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/dijkstracula/go-wortteiler/dictionary"
	"github.com/dijkstracula/go-wortteiler/server"
	"github.com/dijkstracula/go-wortteiler/splitter"
)

const (
	wordPath = "db/de_words.txt"
	prefPath = "db/de_prefixes.txt"
	suffPath = "db/de_suffixes.txt"
)

var (
	logPath = flag.String("logPath", "", "log file path (if unset, use stderr)")
	port    = flag.Int("port", 8080, "The port to listen on")
)

func handleLogPath(logPath string) {
	if logPath != "" {
		f, err := os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			log.Panic("Can't open log file for writing: %v", err)
		}
		log.SetOutput(f)
	}
}

func main() {
	fmt.Println("~~~ go-wortteiler ~~~ by Nathan Taylor <nbtaylor@gmail.com>")

	flag.Parse()

	handleLogPath(*logPath)

	lookup, err := dictionary.FromFiles(wordPath, prefPath, suffPath)
	if err != nil {
		log.Fatalf("Dictionary creation failed: %v", err)
	}

	spl := splitter.Splitter(dictionary.ValidFunc(lookup))

	handlers := server.New(spl)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: handlers,
	}

	log.Printf("Beginning to listen on port %d...\n", *port)
	log.Fatal(server.ListenAndServe())
}
