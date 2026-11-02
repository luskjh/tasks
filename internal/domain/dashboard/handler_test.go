package dashboard

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"sorantics.com/tasks/internal/lib"
	"sorantics.com/tasks/internal/template"
)

func TestOverview(t *testing.T) {
	reg, err := lib.NewRegistry(template.FS)
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Overview(reg)(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("expected response code to be %d, got %d", http.StatusOK, rr.Code)
	}

	for _, expected := range []string{
		"<title>Overview</title>",
		"<h1>Overview</h1>",
		"<nav>",
	} {
		if body := rr.Body.String(); !strings.Contains(body, expected) {
			t.Errorf("expected response body to contain %q, got %q", expected, body)
		}
	}
}
