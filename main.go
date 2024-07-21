package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"tailscale.com/tsnet"
)

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
	r.HandleFunc("/request", server.SongRequest)
	r.HandleFunc("/queue", server.Queue)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./web/")))

	http.Serve(ln, r)
}
