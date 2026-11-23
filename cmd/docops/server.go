package main

import (
	"io"
	"net/http"
	"time"
)

// helloServer embeds the standard HTTP server so docops can serve its
// greeting without pulling in any external dependencies.
type helloServer struct {
	*http.Server
}

func newHelloServer(addr string) *helloServer {
	mux := http.NewServeMux()
	mux.Handle("/", helloHandler{})
	mux.Handle("/hello", helloHandler{})

	return &helloServer{
		Server: &http.Server{
			Addr:              addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
	}
}

// helloHandler writes the greeting as plain text.
type helloHandler struct{}

func (helloHandler) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	_, _ = io.WriteString(w, "Hello World!")
}
