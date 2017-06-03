package server

import (
	"net/http"

	"github.com/dijkstracula/go-wortteiler/dictionary"
	"github.com/gorilla/mux"
)

func splitFunc(d *dictionary.Dictionary) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("todo"))
	}
}

// New produces a new HTTP handler with the appropriate endpoints configured.
func New(d *dictionary.Dictionary) http.Handler {
	r := mux.NewRouter()

	r.NewRoute().
		Methods("POST").
		Path("/split").
		HandlerFunc(splitFunc(d))
	r.NewRoute().
		Methods("GET").
		PathPrefix("/").
		Handler(http.FileServer(http.Dir("public/")))

	return r
}
