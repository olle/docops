package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var docopsBinary string

func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "docops-e2e-")
	if err != nil {
		panic(err)
	}
	docopsBinary = filepath.Join(dir, "docops")
	if runtime.GOOS == "windows" {
		docopsBinary += ".exe"
	}

	repoRoot := filepath.Join("..", "..")
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
