package versionedconfig

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"pgregory.net/rapid"
)

// Validates: Requirements 1.2, 2.3
// Property 1: Error chain preservation
// For any error returned by a wrapping code path (read failure, unmarshal failure,
// decode failure), calling errors.Is(returnedErr, originalCause) SHALL return true,
// confirming the original error is preserved in the chain.

// propTestConfig implements VersionedConfig for property tests.
type propTestConfig struct {
	Kind          string `mapstructure:"kind"`
	SchemaVersion string `mapstructure:"schemaVersion"`
}

func (c *propTestConfig) GetKind() string    { return c.Kind }
func (c *propTestConfig) GetVersion() string { return c.SchemaVersion }

// strictConfig requires a specific typed field that will fail mapstructure decode
// when the YAML provides an incompatible type.
type strictConfig struct {
	Kind          string `mapstructure:"kind"`
	SchemaVersion string `mapstructure:"schemaVersion"`
	Port          int    `mapstructure:"port"`
}

func (c *strictConfig) GetKind() string    { return c.Kind }
func (c *strictConfig) GetVersion() string { return c.SchemaVersion }

func TestProperty_ErrorChainPreservation_ReadConfig(t *testing.T) {
	// "read config" wrapping path: non-existent file triggers os.ErrNotExist
	rapid.Check(t, func(t *rapid.T) {
		// Generate a random directory name that won't exist
		dirName := rapid.StringMatching(`[a-z]{5,10}`).Draw(t, "dirName")
		fileName := rapid.StringMatching(`[a-z]{3,8}\.yaml`).Draw(t, "fileName")
		nonExistentPath := filepath.Join("/tmp", "nonexistent_"+dirName, fileName)

		versions := Versions{
			{SchemaVersion: "v1", Kind: "Test", Factory: func() VersionedConfig { return &propTestConfig{} }},
		}

		_, err := New(nonExistentPath, versions)
		if err == nil {
			t.Fatal("expected error for non-existent file")
		}

		// The original cause (os.ErrNotExist) must be preserved in the error chain
		if !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("errors.Is(err, os.ErrNotExist) = false; err = %v", err)
		}
	})
}

func TestProperty_ErrorChainPreservation_UnmarshalConfig(t *testing.T) {
	// "unmarshal config" wrapping path: invalid YAML triggers a yaml parse error
	rapid.Check(t, func(t *rapid.T) {
		// Generate invalid YAML content that will fail to parse
		invalidContent := rapid.SampledFrom([]string{
			"{not valid yaml: [[[",
			":\n  :\n    - [invalid",
			"{{{{",
			"key: [unclosed",
			"\t\t\t---\n:::invalid:::",
		}).Draw(t, "invalidYAML")

		dir, err := os.MkdirTemp("", "prop_test_*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(dir) }()

		path := filepath.Join(dir, "config.yaml")
		if writeErr := os.WriteFile(path, []byte(invalidContent), 0o644); writeErr != nil {
			t.Fatalf("failed to write test file: %v", writeErr)
		}

		versions := Versions{
			{SchemaVersion: "v1", Kind: "Test", Factory: func() VersionedConfig { return &propTestConfig{} }},
		}

		_, err = New(path, versions)
		if err == nil {
			t.Fatal("expected error for invalid YAML")
		}

		// The error must contain a wrapped cause (the yaml parse error).
		// We verify the chain is preserved by checking Unwrap is non-nil.
		var unwrapper interface{ Unwrap() error }
		if !errors.As(err, &unwrapper) {
			t.Fatalf("error does not implement Unwrap(); err = %v", err)
		}
		if unwrapper.Unwrap() == nil {
			t.Fatalf("Unwrap() returned nil; err = %v", err)
		}
	})
}

func TestProperty_ErrorChainPreservation_ParseConfig(t *testing.T) {
	// "parse config" wrapping path: valid YAML with registered kind/schemaVersion
	// but a field type mismatch that causes mapstructure.Decode to fail.
	rapid.Check(t, func(t *rapid.T) {
		// Generate a non-numeric string that will fail to decode into an int field
		badPort := rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "badPort")

		yamlContent := "kind: Strict\nschemaVersion: v1\nport: " + badPort + "\n"

		dir, err := os.MkdirTemp("", "prop_test_*")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer func() { _ = os.RemoveAll(dir) }()

		path := filepath.Join(dir, "config.yaml")
		if writeErr := os.WriteFile(path, []byte(yamlContent), 0o644); writeErr != nil {
			t.Fatalf("failed to write test file: %v", writeErr)
		}

		versions := Versions{
			{SchemaVersion: "v1", Kind: "Strict", Factory: func() VersionedConfig { return &strictConfig{} }},
		}

		_, err = New(path, versions)
		if err == nil {
			t.Fatal("expected error for type mismatch in mapstructure decode")
		}

		// The error must contain a wrapped cause (the mapstructure decode error).
		// We verify the chain is preserved by checking Unwrap is non-nil.
		var unwrapper interface{ Unwrap() error }
		if !errors.As(err, &unwrapper) {
			t.Fatalf("error does not implement Unwrap(); err = %v", err)
		}
		if unwrapper.Unwrap() == nil {
			t.Fatalf("Unwrap() returned nil; err = %v", err)
		}
	})
}
