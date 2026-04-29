package main

import (
	"fmt"
	lib "lab02_02/lab02_02lib"
	"net"
	"net/http"
)

const A01 = 3

func roothandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "A01: %d\nA02: %t\nA03: %s", A01, A02, lib.A03)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	l, e := net.Listen("tcp", ":4000")
	if e != nil {
		println(e.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", roothandler)
	http.Serve(l, mux)
}
