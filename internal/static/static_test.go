package static

import (
	"errors"
	"io/fs"
	"strings"
	"testing"
	"testing/fstest"
)

type errorFS struct{}

func (errorFS) Open(string) (fs.File, error) {
	return nil, errors.New("open failed")
}

func (errorFS) Sub(string) (fs.FS, error) {
	return nil, errors.New("sub failed")
}

var testFS = fstest.MapFS{
	"dist/foo.txt":  {Data: []byte("baz")},
	"other/bar.txt": {Data: []byte("quux")},
}

func TestDistFS(t *testing.T) {
	distFS, err := DistFS(testFS)
	if err != nil {
		t.Fatalf("DistFS() returned an unexpected error: %v", err)
	}

	data, err := fs.ReadFile(distFS, "foo.txt")
	if err != nil {
		t.Fatalf(`ReadFile() error: %v`, err)
	}
	content := string(data)

	expected := "baz"
	if content != expected {
		t.Errorf("foo.html = %q, want %q", content, expected)
	}

	for _, path := range []string{
		"other/bar.txt",
		"bar.txt",
	} {
		if _, err := fs.ReadFile(distFS, path); err == nil {
			t.Errorf("expected error reading file outside dist (%q), got nil", path)
		}

	}

}

func TestDistFS_SubError(t *testing.T) {
	distFS, err := DistFS(errorFS{})

	if distFS != nil {
		t.Errorf("expected distFS to be nil, got %v", distFS)
	}
	if err == nil {
		t.Fatalf("expected DistFS() to return a nonnil error")
	}

	expected := "sub failed"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}
