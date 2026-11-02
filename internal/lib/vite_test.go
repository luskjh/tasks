package lib

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewVite(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(`{
		"js/dashboard.js": {
			"file": "assets/dashboard.js",
			"isEntry": true,
			"imports": ["_shared.js"],
			"css": ["assets/dashboard.css"]
		},
		"_shared.js": {
			"file": "assets/shared.js",
			"imports": ["_nested.js"],
			"css": ["assets/shared.css"]
		},
		"_nested.js": {
			"file": "assets/nested.js",
			"css": ["assets/nested.css"]
		},
		"js/other.js": {
			"file": "assets/other.js",
			"isEntry": true,
			"css": ["assets/other.css"]
		}
	}`), 0644); err != nil {
		t.Fatalf("error writing test manifest: %v", err)
	}

	vite, err := NewVite(path)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	if l := len(*vite); l != 2 {
		t.Fatalf("expected two entries, got %d", l)
	}

	entry, ok := (*vite)["js/dashboard.js"]
	if !ok {
		t.Fatalf("expected entry map to contain %q", "js/dashboard.js")
	}

	want := `
  <link rel="modulepreload" href="/assets/nested.js" />

  <link rel="modulepreload" href="/assets/shared.js" />


  <link rel="stylesheet" href="/assets/nested.css" />

  <link rel="stylesheet" href="/assets/shared.css" />

  <link rel="stylesheet" href="/assets/dashboard.css" />


  <script type="module" src="/assets/dashboard.js"></script>

`
	if string(entry) != want {
		t.Errorf("expected entry HTML:\n%s\ngot:\n%s", want, entry)
	}

	if got := vite.Entry("js/dashboard.js"); got != entry {
		t.Errorf("expected Entry(%q) to return %q, got %q", "js/dashboard.js", entry, got)
	}

	if got := vite.Entry("missing.js"); got != "" {
		t.Errorf("expected missing entry to be empty, got %q", got)
	}

	if _, ok := (*vite)["_shared.js"]; ok {
		t.Errorf("expected non-entry chunks to be omitted")
	}
}

func TestNewVite_DeduplicatesDependencies(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(`{
		"js/dashboard.js": {
			"file": "assets/dashboard.js",
			"isEntry": true,
			"imports": ["_a.js", "_b.js"],
			"css": ["assets/dashboard.css"]
		},
		"_a.js": {
			"file": "assets/a.js",
			"imports": ["_shared.js"],
			"css": ["assets/a.css"]
		},
		"_b.js": {
			"file": "assets/b.js",
			"imports": ["_shared.js"],
			"css": ["assets/b.css"]
		},
		"_shared.js": {
			"file": "assets/shared.js",
			"imports": ["_a.js"],
			"css": ["assets/shared.css"]
		}
	}`), 0644); err != nil {
		t.Fatalf("error writing test manifest: %v", err)
	}

	vite, err := NewVite(path)
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	entry, ok := (*vite)["js/dashboard.js"]
	if !ok {
		t.Fatalf("expected entry map to contain %q", "js/dashboard.js")
	}

	want := `
  <link rel="modulepreload" href="/assets/shared.js" />

  <link rel="modulepreload" href="/assets/a.js" />

  <link rel="modulepreload" href="/assets/b.js" />


  <link rel="stylesheet" href="/assets/shared.css" />

  <link rel="stylesheet" href="/assets/a.css" />

  <link rel="stylesheet" href="/assets/b.css" />

  <link rel="stylesheet" href="/assets/dashboard.css" />


  <script type="module" src="/assets/dashboard.js"></script>

`
	if string(entry) != want {
		t.Errorf("expected entry HTML:\n%s\ngot:\n%s", want, entry)
	}
}

func TestNewVite_ReadFileError(t *testing.T) {
	vite, err := NewVite(filepath.Join(t.TempDir(), ".vite", "manifest.json"))
	if vite != nil {
		t.Errorf("expected entry map to be nil, got %v", vite)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}
}

func TestNewVite_UnmarshalError(t *testing.T) {
	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(`{`), 0644); err != nil {
		t.Fatalf("error writing test manifest: %v", err)
	}

	vite, err := NewVite(path)
	if vite != nil {
		t.Errorf("expected entry map to be nil, got %v", vite)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}
}
