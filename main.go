package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/netutil"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func logHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h.ServeHTTP(w, r)
		fmt.Printf(
			"%s - - [%s] \"%s %s %s %s\"\n",
			r.RemoteAddr,
			time.Now(),
			r.Method,
			r.URL,
			r.Proto,
			"200",
		)
	})
}

func main() {
	var (
		host     string
		port_str string
		dir      string
	)
	flag.StringVar(&host, "host", "0.0.0.0", "host")
	flag.StringVar(&port_str, "p", "80", "port")
	flag.StringVar(&dir, "d", "./", "serve directory path")
	flag.Parse()
	port, err := strconv.ParseUint(port_str, 10, 16)
	if err != nil {
		log.Fatalln(err)
	}

	mux := http.NewServeMux()
	files := http.FileServer(http.Dir(dir))
	mux.Handle("/", logHandler(http.StripPrefix("/", files)))
	server := &http.Server{
		Handler: mux,
	}

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))

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
