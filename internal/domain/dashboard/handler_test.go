package dashboard

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOverview(t *testing.T) {
	rr := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/", nil)

	Overview(rr, r)

	if rr.Code != http.StatusOK {
		t.Errorf("expected response code to be %d, got %d", http.StatusOK, rr.Code)
	}

	expected := "overview"
	if body := rr.Body.String(); body != expected {
		t.Errorf("expected response body to be %q, got %q", expected, body)
	}
}
