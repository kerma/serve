package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
)

type ResponseWriter struct {
	http.ResponseWriter
	StatusCode int
}

func (r *ResponseWriter) WriteHeader(status int) {
	r.StatusCode = status
	r.ResponseWriter.WriteHeader(status)
}

func handler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := &ResponseWriter{ResponseWriter: w, StatusCode: 200}
		h.ServeHTTP(rw, r)
		log.Printf("%s %s %d\n", r.Method, r.URL, rw.StatusCode)
	}
}

func main() {
	dir := flag.String("d", ".", "directory to serve")
	port := flag.String("p", "9000", "Listen on port")
	flag.Parse()

	fs := http.FileServer(http.Dir(filepath.Dir(*dir)))

	log.Printf("Serving %s on http://[::]:%s\n", *dir, *port)
	log.Fatal(http.ListenAndServe(":"+*port, handler(fs)))
}
