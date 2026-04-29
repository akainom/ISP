package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

func roothandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)
}

func main() {
	file, e := os.OpenFile("03_01.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if e != nil {
		log.Fatal("Failed to open log file:", e)
	}

	log.SetOutput(file)

	l, e := net.Listen("tcp", ":4000")
	if e != nil {
		println(e.Error())
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", roothandler)
	http.Serve(l, mux)
}
