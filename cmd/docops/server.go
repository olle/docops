package main

import (
	"errors"
	"log"
	"net/http"
	"time"
)

type httpServer interface {
	ListenAndServe(addr string, handler http.Handler) error
}

type stdHTTPServer struct{}

func (stdHTTPServer) ListenAndServe(addr string, handler http.Handler) error {
	server := &http.Server{
		Addr:              addr,
		Handler:           handler,
		ReadHeaderTimeout: 5 * time.Second,
	}
	err := server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func handlerFor(writer outputWriter, titles []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		if err := writer.Render(w, "World", titles); err != nil {
			log.Printf("rendering response: %v", err)
		}
	})
}
