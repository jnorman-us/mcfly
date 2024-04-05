package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
	"time"
)

func main() {
	s := &http.Server{
		Addr:           "fly-local-6pn:8080",
		Handler:        HelloHandler{},
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

type HelloHandler struct{}

func (h HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}
