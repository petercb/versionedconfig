package versionedconfig

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetConfigType(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     string
	}{
		{name: "yaml extension", filename: "config.yaml", want: "yaml"},
		{name: "yml extension", filename: "config.yml", want: "yaml"},
		{name: "YAML uppercase", filename: "config.YAML", want: "yaml"},
		{name: "YML uppercase", filename: "config.YML", want: "yaml"},
		{name: "json extension", filename: "config.json", want: "json"},
		{name: "toml extension", filename: "config.toml", want: "toml"},
		{name: "no extension", filename: "config", want: ""},
		{name: "dot only", filename: ".", want: ""},
		{name: "empty string", filename: "", want: ""},
		{name: "path with yaml", filename: "/etc/app/config.yaml", want: "yaml"},
		{name: "path with yml", filename: "/etc/app/config.yml", want: "yaml"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getConfigType(tt.filename)
			if got != tt.want {
				t.Errorf("getConfigType(%q) = %q, want %q", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsURL(t *testing.T) {
	tests := []struct {
		name string
		s    string
		want bool
	}{
		{name: "https url", s: "https://example.com/config.yaml", want: true},
		{name: "http url", s: "http://example.com/config.yaml", want: true},
		{name: "file path", s: "/etc/config.yaml", want: false},
		{name: "relative path", s: "config.yaml", want: false},
		{name: "ftp url", s: "ftp://example.com/config.yaml", want: false},
		{name: "empty string", s: "", want: false},
		{name: "https no path", s: "https://example.com", want: true},
		{name: "http no path", s: "http://example.com", want: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isURL(tt.s)
			if got != tt.want {
				t.Errorf("isURL(%q) = %v, want %v", tt.s, got, tt.want)
			}
		})
	}
}

func TestReadConfiguration(t *testing.T) {
	t.Run("empty filename returns error", func(t *testing.T) {
		_, err := readConfiguration("")
		if err == nil {
			t.Fatal("expected error for empty filename")
		}
	})

	t.Run("nonexistent file returns error", func(t *testing.T) {
		_, err := readConfiguration("/nonexistent/path/config.yaml")
		if err == nil {
			t.Fatal("expected error for nonexistent file")
		}
	})

	t.Run("reads existing file", func(t *testing.T) {
		dir := t.TempDir()
		path := filepath.Join(dir, "test.yaml")
		content := []byte("kind: Config\nschemaVersion: v1\n")
		if err := os.WriteFile(path, content, 0o644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}

		got, err := readConfiguration(path)
		if err != nil {
			t.Fatalf("readConfiguration(%q) error: %v", path, err)
		}
		if string(got) != string(content) {
			t.Errorf("readConfiguration(%q) = %q, want %q", path, got, content)
		}
	})
}
