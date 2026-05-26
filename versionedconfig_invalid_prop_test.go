package versionedconfig

import (
	"os"
	"path/filepath"
	"testing"

	"pgregory.net/rapid"
)

// **Validates: Requirements 4.1**

// invalidTestVersions returns a Versions slice used for invalid input testing.
func invalidTestVersions() Versions {
	return Versions{
		{SchemaVersion: "v1", Kind: "App", Factory: func() VersionedConfig { return &testConfig{} }},
		{SchemaVersion: "v2", Kind: "App", Factory: func() VersionedConfig { return &testConfig{} }},
	}
}

// writeInvalidYAML writes content to a temp file with the given extension and returns the path.
func writeInvalidYAML(t *rapid.T, content string, ext string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "versionedconfig-invalid-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	path := filepath.Join(dir, "config"+ext)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}
	return path
}

func TestProperty_InvalidInput_EmptyFilename(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate arbitrary Versions slices — the filename is always empty
		_, err := New("", invalidTestVersions())
		if err == nil {
			t.Fatal("expected non-nil error for empty filename")
		}
	})
}

func TestProperty_InvalidInput_NonExistentFile(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random directory-like paths that don't exist
		segments := rapid.SliceOfN(rapid.StringMatching(`[a-z]{3,8}`), 2, 5).Draw(t, "pathSegments")
		filename := rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "filename")
		path := "/" + filepath.Join(segments...) + "/" + filename + ".yaml"

		_, err := New(path, invalidTestVersions())
		if err == nil {
			t.Fatalf("expected non-nil error for non-existent file: %s", path)
		}
	})
}

func TestProperty_InvalidInput_MissingKind(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with schemaVersion but no kind field
		version := rapid.SampledFrom([]string{"v1", "v2", "v3", "v99"}).Draw(t, "version")
		name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9]{0,15}`).Draw(t, "name")

		content := "schemaVersion: " + version + "\nname: " + name + "\n"
		path := writeInvalidYAML(t, content, ".yaml")

		_, err := New(path, invalidTestVersions())
		if err == nil {
			t.Fatal("expected non-nil error for YAML missing kind field")
		}
	})
}

func TestProperty_InvalidInput_MissingSchemaVersion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with kind but no schemaVersion field
		kind := rapid.SampledFrom([]string{"App", "Service", "Config", "Unknown"}).Draw(t, "kind")
		name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9]{0,15}`).Draw(t, "name")

		content := "kind: " + kind + "\nname: " + name + "\n"
		path := writeInvalidYAML(t, content, ".yaml")

		_, err := New(path, invalidTestVersions())
		if err == nil {
			t.Fatal("expected non-nil error for YAML missing schemaVersion field")
		}
	})
}

func TestProperty_InvalidInput_UnknownKindVersion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with kind/schemaVersion combos that are NOT registered
		kind := rapid.SampledFrom([]string{"Unknown", "Service", "Database", "Widget"}).Draw(t, "kind")
		version := rapid.SampledFrom([]string{"v99", "v100", "alpha", "beta"}).Draw(t, "version")

		content := "kind: " + kind + "\nschemaVersion: " + version + "\nname: test\n"
		path := writeInvalidYAML(t, content, ".yaml")

		_, err := New(path, invalidTestVersions())
		if err == nil {
			t.Fatalf("expected non-nil error for unknown kind/version combo: %s/%s", kind, version)
		}
	})
}

func TestProperty_InvalidInput_NonYAMLExtension(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate files with non-YAML extensions
		ext := rapid.SampledFrom([]string{".json", ".toml", ".txt", ".xml", ".ini", ".cfg"}).Draw(t, "ext")
		content := "kind: App\nschemaVersion: v1\nname: test\n"
		path := writeInvalidYAML(t, content, ext)

		_, err := New(path, invalidTestVersions())
		if err == nil {
			t.Fatalf("expected non-nil error for non-YAML extension: %s", ext)
		}
	})
}
