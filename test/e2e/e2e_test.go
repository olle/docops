package e2e_test

import (
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

var (
	docopsBinary string
	repoRoot     = filepath.Join("..", "..")
)

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "docops-e2e-")
	if err != nil {
		panic(err)
	}
	docopsBinary = filepath.Join(dir, "docops")
	if runtime.GOOS == "windows" {
		docopsBinary += ".exe"
	}

	goBin := filepath.Join(runtime.GOROOT(), "bin", "go")
	build := exec.Command(goBin, "build", "-o", docopsBinary, "./cmd/docops")
	build.Dir = repoRoot
	if output, err := build.CombinedOutput(); err != nil {
		panic("building docops: " + err.Error() + "\n" + string(output))
	}

	code := m.Run()
	if err := os.RemoveAll(dir); err != nil {
		panic(err)
	}
	os.Exit(code)
}

func TestHelloNameResolvesName(t *testing.T) {
	output, err := exec.Command(docopsBinary, "-name", "Ada").CombinedOutput()
	if err != nil {
		t.Fatalf("running docops: %v\n%s", err, output)
	}

	if got, want := strings.TrimSpace(string(output)), "Hello Ada!"; got != want {
		t.Errorf("greeting = %q, want %q", got, want)
	}
}

func TestServerModeServesHelloWorld(t *testing.T) {
	addr := freeAddr(t)
	cmd := exec.Command(docopsBinary, "--server", "--addr", addr)
	cmd.Dir = repoRoot
	cmd.Stdout = io.Discard
	cmd.Stderr = io.Discard
	if err := cmd.Start(); err != nil {
		t.Fatalf("starting docops server: %v", err)
	}
	defer func() {
		_ = cmd.Process.Kill()
		_ = cmd.Wait()
	}()

	url := "http://" + addr + "/hello"
	var (
		resp *http.Response
		err  error
	)
	deadline := time.Now().Add(5 * time.Second)
	for {
		resp, err = http.Get(url)
		if err == nil {
			break
		}
		if time.Now().After(deadline) {
			t.Fatalf("server at %s never became ready: %v", addr, err)
		}
		time.Sleep(25 * time.Millisecond)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("reading response: %v", err)
	}
	if got, want := resp.Header.Get("Content-Type"), "text/html; charset=utf-8"; got != want {
		t.Errorf("Content-Type = %q, want %q", got, want)
	}
	got := string(body)
	if !strings.Contains(got, "<!doctype html>") ||
		!strings.Contains(got, "<h1>Hello World!</h1>") ||
		!strings.Contains(got, "<li>Hello World Events</li>") {
		t.Errorf("response body = %q, want semantic HTML with greeting and AsyncAPI title", got)
	}
	if !strings.Contains(got, "<style>") || !strings.Contains(got, ":not(pre) > code") {
		t.Errorf("response body = %q, want inlined CSS", got)
	}
	if strings.Contains(got, `<link rel="stylesheet"`) {
		t.Errorf("response body = %q, want no external stylesheet", got)
	}
}

func freeAddr(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("finding a free port: %v", err)
	}
	addr := listener.Addr().String()
	if err := listener.Close(); err != nil {
		t.Fatalf("releasing free port: %v", err)
	}
	return addr
}
