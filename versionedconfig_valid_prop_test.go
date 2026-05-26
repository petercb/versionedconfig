package versionedconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"pgregory.net/rapid"
)

// propValidConfig implements VersionedConfig for property testing valid inputs.
type propValidConfig struct {
	Kind          string `mapstructure:"kind"`
	SchemaVersion string `mapstructure:"schemaVersion"`
	Name          string `mapstructure:"name"`
	Value         int    `mapstructure:"value"`
	Enabled       bool   `mapstructure:"enabled"`
}

func (c *propValidConfig) GetKind() string    { return c.Kind }
func (c *propValidConfig) GetVersion() string { return c.SchemaVersion }

// TestProperty4_ValidInputSuccessPreservation verifies that for any valid YAML
// file containing a registered kind and schemaVersion with well-formed fields,
// the library returns a non-nil VersionedConfig and a nil error.
//
// **Validates: Requirements 4.2**
func TestProperty4_ValidInputSuccessPreservation(t *testing.T) {
	const kind = "PropValid"
	const schemaVersion = "v1"

	versions := Versions{
		{
			SchemaVersion: schemaVersion,
			Kind:          kind,
			Factory:       func() VersionedConfig { return &propValidConfig{} },
		},
	}

	rapid.Check(t, func(rt *rapid.T) {
		// Generate varied but valid field values.
		// Use a pattern that avoids YAML 1.1 boolean literals (y, n, yes, no, on, off)
		// by requiring at least 2 characters starting with a letter.
		name := rapid.StringMatching(`[a-zA-Z][a-zA-Z0-9_-]{1,20}`).Draw(rt, "name")
		value := rapid.IntRange(0, 10000).Draw(rt, "value")
		enabled := rapid.Bool().Draw(rt, "enabled")

		// Build valid YAML content with required kind and schemaVersion fields.
		// Quote the name value to prevent YAML interpreting it as a non-string type.
		yamlContent := fmt.Sprintf(
			"kind: %s\nschemaVersion: %s\nname: \"%s\"\nvalue: %d\nenabled: %v\n",
			kind, schemaVersion, name, value, enabled,
		)

		// Write YAML to a temp file
		dir, err := os.MkdirTemp("", "prop-valid-*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(dir) }()

		path := filepath.Join(dir, "config.yaml")
		if writeErr := os.WriteFile(path, []byte(yamlContent), 0o644); writeErr != nil {
			t.Fatalf("failed to write temp YAML: %v", writeErr)
		}

		// Call New() with matching versions
		result, err := New(path, versions)
		// Assert: nil error and non-nil result
		if err != nil {
			rt.Fatalf("New() returned unexpected error for valid input: %v\nYAML:\n%s", err, yamlContent)
		}
		if result == nil {
			rt.Fatal("New() returned nil VersionedConfig for valid input")
		}

		// Verify the result has correct kind and version
		if result.GetKind() != kind {
			rt.Errorf("GetKind() = %q, want %q", result.GetKind(), kind)
		}
		if result.GetVersion() != schemaVersion {
			rt.Errorf("GetVersion() = %q, want %q", result.GetVersion(), schemaVersion)
		}
	})
}
