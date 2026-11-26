package main

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olle/docops/internal/view"
)

type fakeSnapshotSource struct {
	snapshot view.Snapshot
	err      error
}

func (f fakeSnapshotSource) Snapshot(context.Context) (view.Snapshot, error) {
	return f.snapshot, f.err
}

type recordingPresenter struct {
	snapshot view.Snapshot
}

func (p *recordingPresenter) Render(w io.Writer, snapshot view.Snapshot) error {
	p.snapshot = snapshot
	_, err := io.WriteString(w, "rendered:"+snapshot.Greeting)
	return err
}

func (*recordingPresenter) ContentType() string {
	return "text/fake; charset=utf-8"
}

type recordingScreen struct {
	ran bool
}

func (s *recordingScreen) Run(context.Context) error {
	s.ran = true
	return nil
}

type recordingHTTPServer struct {
	addr    string
	handler http.Handler
}

func (s *recordingHTTPServer) ListenAndServe(addr string, handler http.Handler) error {
	s.addr = addr
	s.handler = handler
	return nil
}

func TestGreetCommandUsesOnlyInformationAndPresenter(t *testing.T) {
	var out bytes.Buffer
	source := fakeSnapshotSource{snapshot: view.Snapshot{Greeting: "Hello Ada!"}}
	presenter := &recordingPresenter{}
	if err := (greetCommand{info: source, presenter: presenter, out: &out}).Run(); err != nil {
		t.Fatal(err)
	}
	if out.String() != "rendered:Hello Ada!" {
		t.Errorf("output = %q, want %q", out.String(), "rendered:Hello Ada!")
	}
	if presenter.snapshot.Greeting != source.snapshot.Greeting {
		t.Errorf("presented %q, want %q", presenter.snapshot.Greeting, source.snapshot.Greeting)
	}
}

func TestStaticCommandWritesAndUsesPresenter(t *testing.T) {
	var out bytes.Buffer
	source := fakeSnapshotSource{snapshot: view.Snapshot{Greeting: "Hello World!"}}
	presenter := &recordingPresenter{}
	if err := (staticCommand{info: source, presenter: presenter, out: &out}).Run(); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), "rendered:Hello World!"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestTUICommandUsesOnlyScreen(t *testing.T) {
	screen := &recordingScreen{}
	if err := (tuiCommand{screen: screen}).Run(); err != nil {
		t.Fatal(err)
	}
	if !screen.ran {
		t.Error("screen was not run")
	}
}

func TestServeCommandUsesOnlyInformationAndPresenter(t *testing.T) {
	source := fakeSnapshotSource{snapshot: view.Snapshot{Greeting: "Hello Ada!"}}
	presenter := &recordingPresenter{}
	server := &recordingHTTPServer{}
	if err := (serveCommand{addr: "127.0.0.1:0", info: source, presenter: presenter, server: server}).Run(); err != nil {
		t.Fatal(err)
	}
	if server.addr != "127.0.0.1:0" || server.handler == nil {
		t.Fatalf("server addr = %q, handler = %v", server.addr, server.handler)
	}

	recorder := httptest.NewRecorder()
	server.handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))
	if got := recorder.Header().Get("Content-Type"); got != "text/fake; charset=utf-8" {
		t.Errorf("Content-Type = %q, want %q", got, "text/fake; charset=utf-8")
	}
	if got := recorder.Body.String(); got != "rendered:Hello Ada!" {
		t.Errorf("body = %q, want %q", got, "rendered:Hello Ada!")
	}
}

func TestDebugCommandDumpsSnapshot(t *testing.T) {
	var out bytes.Buffer
	snapshot := view.Snapshot{Documents: []view.Document{{KindLabel: "AsyncAPI", DisplayName: "Events API"}}}
	command := debugCommand{source: fakeSnapshotSource{snapshot: snapshot}, out: &out}
	if err := command.Run(); err != nil {
		t.Fatal(err)
	}
	if got, want := out.String(), view.Debug(snapshot); got != want {
		t.Errorf("debug output = %q, want %q", got, want)
	}
}

func TestDebugCommandPropagatesSourceError(t *testing.T) {
	wantErr := errors.New("snapshot unavailable")
	command := debugCommand{source: fakeSnapshotSource{err: wantErr}, out: io.Discard}
	if err := command.Run(); !errors.Is(err, wantErr) {
		t.Errorf("error = %v, want %v", err, wantErr)
	}
}
