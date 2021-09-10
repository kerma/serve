package main

import (
	"flag"
	"log"
	"net/http"
	"os"
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

func main() {

	port := flag.String("p", "9000", "Listen on port")
	flag.Parse()

	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	path := filepath.Dir(ex)

	fs := http.FileServer(http.Dir(path))
	handler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &ResponseWriter{ResponseWriter: w, StatusCode: 200}
			fs.ServeHTTP(rw, r)
			log.Printf("%s %s %d\n", r.Method, r.URL, rw.StatusCode)
		})
	}

	log.Printf("Serving http://[::]:%s\n", *port)
	log.Fatal(http.ListenAndServe(":"+*port, handler()))
}
