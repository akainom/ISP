package main

import (
	"fmt"
	lib "lab02_01/go02_01lib"
	"net"
	"net/http"
)

const c01 = 3.14

func roothandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		fmt.Fprintf(w, "c01: %f\nc02: %e\nc03: %f", c01, c02, lib.C03)
		fmt.Printf("%f\n", c01)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func main() {
	l, err := net.Listen("tcp", ":3000")
	if err != nil {
		print(err.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", roothandler)

	http.Serve(l, mux)
}
