package main

import (
	"golang.org/x/net/netutil"
	"log"
	"net"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("./"))
	mux.Handle("/", http.StripPrefix("/", files))
	server := &http.Server{
		Handler: mux,
	}

	ln, err := net.Listen("tcp", "0.0.0.0:8080")

	if err != nil {
		log.Fatalln(err)
	}

	limit_ln := netutil.LimitListener(ln, 100)
	defer limit_ln.Close()

	err = server.Serve(limit_ln)

	if err != nil {
		log.Fatalln(err)
	}
}
