package lib

import (
	"os"
	"strings"
	"testing"
)

func TestNewConfig(t *testing.T) {
	path := t.TempDir() + "config.json"
	if err := os.WriteFile(path, []byte(`{"port":8000}`), 0644); err != nil {
		t.Fatalf("unable to write test config: %v", err)
	}

	config, err := NewConfig(path)
	if config == nil {
		t.Errorf("expected config to be nonnil, got %v", config)
	}
	if err != nil {
		t.Fatalf("expected error to be nil, got %v", err)
	}

	port := uint16(8000)
	if config.Port != port {
		t.Errorf("expected port to be %d, got %d", port, config.Port)
	}
}

func TestNewConfig_ReadFileError(t *testing.T) {
	path := t.TempDir() + "config.json"

	config, err := NewConfig(path)
	if config != nil {
		t.Errorf("expected config to be nil, got %v", config)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}

	expected := "config.json: no such file or directory"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}

func TestNewConfig_UnmarshallError(t *testing.T) {
	path := t.TempDir() + "config.json"
	if err := os.WriteFile(path, []byte("{"), 0644); err != nil {
		t.Fatalf("unable to write test config: %v", err)
	}

	config, err := NewConfig(path)
	if config != nil {
		t.Errorf("expected config to be nil, got %v", config)
	}
	if err == nil {
		t.Fatalf("expected error to be nonnil, got nil")
	}

	expected := "unexpected end of JSON input"
	if msg := err.Error(); !strings.Contains(msg, expected) {
		t.Fatalf("expected error message to contain %q, got %q", expected, msg)
	}
}
