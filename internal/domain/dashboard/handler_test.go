package dashboard

import (
	htemplate "html/template"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"sorantics.com/tasks/internal/lib"
	"sorantics.com/tasks/internal/template"
)

func testReg(t *testing.T) lib.Registry {
	t.Helper()

	path := filepath.Join(t.TempDir(), "manifest.json")
	if err := os.WriteFile(path, []byte(`{
		"js/dashboard.js": {
			"file": "assets/dashboard.js",
			"isEntry": true,
			"css": ["assets/dashboard.css"]
		}
	}`), 0644); err != nil {
		t.Fatalf("error writing test manifest: %v", err)
	}

	vite, err := lib.NewVite(path)
	if err != nil {
		t.Fatalf("error creating vite: %v", err)
	}

	funcMap := htemplate.FuncMap{"vite_entry": vite.Entry}
	reg, err := lib.NewRegistry(template.FS, funcMap)
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}

	return reg
}

func TestOverview(t *testing.T) {
	reg := testReg(t)
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Overview(reg)(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("expected response code to be %d, got %d", http.StatusOK, rr.Code)
	}

	for _, expected := range []string{
		"<title>Overview</title>",
		"<h1>Overview</h1>",
	} {
		if body := rr.Body.String(); !strings.Contains(body, expected) {
			t.Errorf("expected response body to contain %q, got %q", expected, body)
		}
	}
}
