package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tailscale.com/tsnet"
)

func noCacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Tell the browser not to cache results
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
		next.ServeHTTP(w, r)
	})
}

func main() {
	srv := new(tsnet.Server)
	srv.Hostname = "djvivo"
	srv.Dir = "ts"

	ln, err := srv.ListenFunnel("tcp", ":8443")
	if err != nil {
		log.Fatal(err)
	}

	server := NewServer("./data")
	http.HandleFunc("/request", server.SongRequest)
	http.HandleFunc("/queue", server.Queue)

	r := mux.NewRouter()
	r.Use(noCacheMiddleware)
	r.HandleFunc("/request", server.SongRequest)
	r.HandleFunc("/queue", server.Queue)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	http.Serve(ln, r)
}
