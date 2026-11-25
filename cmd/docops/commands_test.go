package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

type fakeInformation struct {
	msg message
}

func (f fakeInformation) Message() (message, error) {
	return f.msg, nil
}

type recordingPresenter struct {
	msg message
}

func (p *recordingPresenter) Render(w io.Writer, msg message) error {
	p.msg = msg
	_, err := io.WriteString(w, "rendered:"+msg.greeting)
	return err
}

func (*recordingPresenter) ContentType() string {
	return "text/fake; charset=utf-8"
}

type recordingScreen struct {
	msg message
}

func (s *recordingScreen) Present(msg message) error {
	s.msg = msg
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
	info := fakeInformation{msg: newMessage("Ada", nil)}
	presenter := &recordingPresenter{}
	if err := (greetCommand{info: info, presenter: presenter, out: &out}).Run(); err != nil {
		t.Fatal(err)
	}
	if out.String() != "rendered:Hello Ada!" {
		t.Errorf("output = %q, want %q", out.String(), "rendered:Hello Ada!")
	}
	if presenter.msg.greeting != info.msg.greeting {
		t.Errorf("presented %q, want %q", presenter.msg.greeting, info.msg.greeting)
	}
}

func TestTUICommandUsesOnlyInformationAndPresenter(t *testing.T) {
	info := fakeInformation{msg: newMessage("Ada", []string{"Events API"})}
	screen := &recordingScreen{}
	if err := (tuiCommand{info: info, presenter: screen}).Run(); err != nil {
		t.Fatal(err)
	}
	if screen.msg.greeting != info.msg.greeting || len(screen.msg.titles) != 1 {
		t.Errorf("presented %+v, want %+v", screen.msg, info.msg)
	}
}

func TestServeCommandUsesOnlyInformationAndPresenter(t *testing.T) {
	info := fakeInformation{msg: newMessage("Ada", nil)}
	presenter := &recordingPresenter{}
	server := &recordingHTTPServer{}
	if err := (serveCommand{addr: "127.0.0.1:0", info: info, presenter: presenter, server: server}).Run(); err != nil {
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
