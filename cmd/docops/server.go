package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/olle/docops/internal/app"
	"github.com/olle/docops/internal/view"
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

func handlerFor(presenter app.WebPresenter, snapshot view.Snapshot) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", presenter.ContentType())
		if err := presenter.Render(w, snapshot); err != nil {
			log.Printf("rendering response: %v", err)
		}
	})
}
