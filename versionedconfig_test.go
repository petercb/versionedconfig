package versionedconfig

import (
	"os"
	"path/filepath"
	"testing"
)

// testConfig implements VersionedConfig for integration tests.
type testConfig struct {
	Kind          string `mapstructure:"kind"`
	SchemaVersion string `mapstructure:"schemaVersion"`
	Name          string `mapstructure:"name"`
	Count         int    `mapstructure:"count"`
}

func (c *testConfig) GetKind() string    { return c.Kind }
func (c *testConfig) GetVersion() string { return c.SchemaVersion }

func newTestConfig() VersionedConfig { return &testConfig{} }

func testVersions() Versions {
	return Versions{
		{SchemaVersion: "v1", Kind: "App", Factory: newTestConfig},
		{SchemaVersion: "v2", Kind: "App", Factory: newTestConfig},
	}
}

func writeYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test yaml: %v", err)
	}
	return path
}

func TestNew_Success(t *testing.T) {
	yaml := `kind: App
schemaVersion: v1
name: myapp
count: 42
`
	path := writeYAML(t, yaml)

	cfg, err := New(path, testVersions())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	tc, ok := cfg.(*testConfig)
	if !ok {
		t.Fatalf("expected *testConfig, got %T", cfg)
	}

	if tc.Kind != "App" {
		t.Errorf("Kind = %q, want %q", tc.Kind, "App")
	}
	if tc.SchemaVersion != "v1" {
		t.Errorf("SchemaVersion = %q, want %q", tc.SchemaVersion, "v1")
	}
	if tc.Name != "myapp" {
		t.Errorf("Name = %q, want %q", tc.Name, "myapp")
	}
	if tc.Count != 42 {
		t.Errorf("Count = %d, want %d", tc.Count, 42)
	}
}

func TestNew_V2(t *testing.T) {
	yaml := `kind: App
schemaVersion: v2
name: v2app
count: 7
`
	path := writeYAML(t, yaml)

	cfg, err := New(path, testVersions())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	tc := cfg.(*testConfig)
	if tc.SchemaVersion != "v2" {
		t.Errorf("SchemaVersion = %q, want %q", tc.SchemaVersion, "v2")
	}
	if tc.Name != "v2app" {
		t.Errorf("Name = %q, want %q", tc.Name, "v2app")
	}
}

func TestNew_MissingKind(t *testing.T) {
	yaml := `schemaVersion: v1
name: myapp
`
	path := writeYAML(t, yaml)

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for missing kind")
	}
}

func TestNew_MissingSchemaVersion(t *testing.T) {
	yaml := `kind: App
name: myapp
`
	path := writeYAML(t, yaml)

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for missing schemaVersion")
	}
}

func TestNew_UnknownSchemaVersion(t *testing.T) {
	yaml := `kind: App
schemaVersion: v99
name: myapp
`
	path := writeYAML(t, yaml)

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for unknown schema version")
	}
}

func TestNew_UnknownKind(t *testing.T) {
	yaml := `kind: Unknown
schemaVersion: v1
name: myapp
`
	path := writeYAML(t, yaml)

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for unknown kind")
	}
}

func TestNew_InvalidYAML(t *testing.T) {
	content := `{not valid yaml: [[[`
	path := writeYAML(t, content)

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for invalid YAML")
	}
}

func TestNew_EmptyFilename(t *testing.T) {
	_, err := New("", testVersions())
	if err == nil {
		t.Fatal("expected error for empty filename")
	}
}

func TestNew_NonexistentFile(t *testing.T) {
	_, err := New("/nonexistent/path/config.yaml", testVersions())
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestNew_UnsupportedFileType(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := []byte(`{"kind": "App", "schemaVersion": "v1"}`)
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	_, err := New(path, testVersions())
	if err == nil {
		t.Fatal("expected error for unsupported file type")
	}
}

func TestNew_YMLExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yml")
	content := []byte("kind: App\nschemaVersion: v1\nname: ymltest\n")
	if err := os.WriteFile(path, content, 0o644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	cfg, err := New(path, testVersions())
	if err != nil {
		t.Fatalf("New() error: %v", err)
	}

	tc := cfg.(*testConfig)
	if tc.Name != "ymltest" {
		t.Errorf("Name = %q, want %q", tc.Name, "ymltest")
	}
}
