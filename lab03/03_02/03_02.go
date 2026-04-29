package main

import (
	"03_02/P03_02"
	"net"
	"net/http"
)

var s = &P03_02.Stats{}

func statshandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.IncGet()
	case "POST":
		s.IncPost()
	}
}

func genstats(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(s.GenStr()))
}

func main() {
	l, e := net.Listen("tcp", ":4000")
	if e != nil {
		println(e.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /S", statshandler)
	mux.HandleFunc("POST /S", statshandler)
	mux.HandleFunc("GET /G", genstats)

	http.Serve(l, mux)
}
