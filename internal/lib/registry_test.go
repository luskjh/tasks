package lib

import (
	"errors"
	"html/template"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

type errorFS struct{}

func (errorFS) Open(name string) (fs.File, error) {
	return nil, errors.New("open failed")
}

type unreadableFS struct{ fstest.MapFS }

func (unreadableFS) Open(name string) (fs.File, error) {
	return nil, fs.ErrNotExist
}

func (unreadableFS) ReadFile(name string) ([]byte, error) {
	return nil, fs.ErrNotExist
}

var (
	testFS = fstest.MapFS{
		"layout/dashboard.html": {
			Data: []byte(`<html>{{ template "shared/head.html" . }}<body>{{ template "content" . }}</body></html>`),
		},
		"page/dashboard/overview.html": {
			Data: []byte(`{{ define "title" }}Overview{{ end }}{{ define "content" }}<h1>Overview</h1><p>{{ hello "fido" }}</p>{{ end }}`),
		},
		"shared/head.html": {
			Data: []byte(`<head><title>{{ template "title" . }}</title></head>`),
		},
		"ignored.txt": {
			Data: []byte("not a template"),
		},
	}
	testFuncs = template.FuncMap{
		"hello": func(name string) string { return "hello, " + name },
	}
)

func TestNewRegistry_WalkError(t *testing.T) {
	registry, err := NewRegistry(errorFS{}, nil)
	if registry != nil {
		t.Errorf("expected registry to be nil, got %v", registry)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}

	expected := "open failed"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}

func TestNewRegistry_BaseParseError(t *testing.T) {
	fsys := fstest.MapFS{
		"layout/dashboard.html": {Data: []byte(`{{`)},
	}

	registry, err := NewRegistry(fsys, nil)
	if registry != nil {
		t.Errorf("expected registry to be nil, got %v", registry)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}

	expected := "unclosed action"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}

func TestNewRegistry_PageMissingError(t *testing.T) {
	fsys := unreadableFS{
		MapFS: fstest.MapFS{
			"page/dashboard/missing.html": {Data: []byte(`{{ define "content" }}Missing{{ end }}`)},
		},
	}

	registry, err := NewRegistry(fsys, nil)
	if registry != nil {
		t.Errorf("expected registry to be nil, got %v", registry)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}

	expected := "file does not exist"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}

func TestRegistry_RenderWithLayout(t *testing.T) {
	registry, err := NewRegistry(testFS, testFuncs)
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}
	rr := httptest.NewRecorder()

	registry.RenderWithLayout(rr, "dashboard/overview.html", nil, "layout/dashboard.html")

	if rr.Code != http.StatusOK {
		t.Errorf("expected response code to be %d, got %d", http.StatusOK, rr.Code)
	}

	body := rr.Body.String()
	for _, expected := range []string{
		"<title>Overview</title>",
		"<h1>Overview</h1>",
		"<p>hello, fido</p>",
	} {
		if !strings.Contains(body, expected) {
			t.Errorf("expected response body to contain %q, got %q", expected, body)
		}
	}
}

func TestRegistry_RenderWithLayout_PageNotFound(t *testing.T) {
	registry, err := NewRegistry(fstest.MapFS{}, nil)
	if err != nil {
		t.Fatalf("error creating registry: %v", err)
	}
	rr := httptest.NewRecorder()

	registry.RenderWithLayout(rr, "dashboard/missing.html", nil, "layout/dashboard.html")

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected response code to be %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	expected := "internal server error"
	if body := rr.Body.String(); !strings.Contains(body, expected) {
		t.Errorf("expected response body to contain %q, got %q", expected, body)
	}
}

func TestRegistry_RenderWithLayout_ExecuteError(t *testing.T) {
	registry := Registry{
		"dashboard/overview.html": template.Must(template.New("page").Parse(`{{ define "layout/dashboard.html" }}{{ template "missing.html" . }}{{ end }}`)),
	}
	rr := httptest.NewRecorder()

	registry.RenderWithLayout(rr, "dashboard/overview.html", nil, "dashboard.html")

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected response code to be %d, got %d", http.StatusInternalServerError, rr.Code)
	}

	expected := "internal server error"
	if body := rr.Body.String(); !strings.Contains(body, expected) {
		t.Errorf("expected response body to contain %q, got %q", expected, body)
	}
}
