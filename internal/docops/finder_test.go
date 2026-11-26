package docops

import (
	"os"
	"path/filepath"
	"slices"
	"testing"
)

func TestDocumentGlobFinder(t *testing.T) {
	root := t.TempDir()
	for _, dir := range []string{".hidden", "bin"} {
		if err := os.Mkdir(filepath.Join(root, dir), 0o755); err != nil {
			t.Fatal(err)
		}
	}

	write := func(path string) {
		t.Helper()
		if err := os.WriteFile(path, []byte("asyncapi: 3.0.0\n"), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for _, path := range []string{
		filepath.Join(root, "data.json"),
		filepath.Join(root, "notes.yml"),
		filepath.Join(root, "spec.yaml"),
		filepath.Join(root, "readme.md"),
		filepath.Join(root, ".hidden", "spec.yaml"),
		filepath.Join(root, "bin", "spec.yaml"),
	} {
		write(path)
	}

	paths, err := (GlobFinder{}).Find(root)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		filepath.Join(root, "data.json"),
		filepath.Join(root, "notes.yml"),
		filepath.Join(root, "spec.yaml"),
	}
	if !slices.Equal(paths, want) {
		t.Errorf("paths = %q, want %q", paths, want)
	}
}
