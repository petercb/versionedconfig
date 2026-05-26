package versionedconfig

import (
	"os"
	"path/filepath"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// **Validates: Requirements 2.1**
// Property 2: Error message format conventions
// For any error returned by the library, the error message string SHALL start
// with a lowercase letter and SHALL NOT end with punctuation (period,
// exclamation, or question mark).

func assertErrorFormat(t *rapid.T, err error) {
	t.Helper()
	if err == nil {
		return
	}
	msg := err.Error()
	if len(msg) == 0 {
		t.Fatal("error message is empty")
	}
	first := rune(msg[0])
	if unicode.IsLetter(first) && unicode.IsUpper(first) {
		t.Fatalf("error message starts with uppercase: %q", msg)
	}
	last := msg[len(msg)-1]
	if last == '.' || last == '!' || last == '?' {
		t.Fatalf("error message ends with punctuation %q: %q", string(last), msg)
	}
}

func writeTempYAML(t *rapid.T, content, ext string) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "proptest-*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	path := filepath.Join(dir, "config."+ext)
	if writeErr := os.WriteFile(path, []byte(content), 0o644); writeErr != nil {
		t.Fatalf("write temp file: %v", writeErr)
	}
	return path
}

func TestProperty_ErrorFormat_EmptyFilename(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Empty filename always triggers "filename not specified"
		_, err := New("", testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_NonexistentFile(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate random directory/file names that won't exist
		name := rapid.StringMatching(`[a-z]{3,10}`).Draw(t, "name")
		filename := filepath.Join("/tmp/nonexistent", name, "config.yaml")
		_, err := New(filename, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_InvalidYAML(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate varied invalid YAML content
		badContent := rapid.SampledFrom([]string{
			"{{{not yaml: [[[",
			":::",
			"\t\t\t---\n\t\t\t- :[",
			"key: [unclosed",
			"{unmatched",
		}).Draw(t, "badContent")

		path := writeTempYAML(t, badContent, "yaml")

		_, err := New(path, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_MissingKind(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with schemaVersion but no kind field
		version := rapid.StringMatching(`v[0-9]{1,3}`).Draw(t, "version")
		content := "schemaVersion: " + version + "\nname: test\n"

		path := writeTempYAML(t, content, "yaml")

		_, err := New(path, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_MissingSchemaVersion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with kind but no schemaVersion field
		kind := rapid.StringMatching(`[A-Z][a-z]{2,8}`).Draw(t, "kind")
		content := "kind: " + kind + "\nname: test\n"

		path := writeTempYAML(t, content, "yaml")

		_, err := New(path, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_UnknownSchemaVersion(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML with kind and schemaVersion that don't match any registered version
		kind := rapid.StringMatching(`[A-Z][a-z]{2,8}`).Draw(t, "kind")
		version := rapid.StringMatching(`v[0-9]{2,4}`).Draw(t, "version")
		content := "kind: " + kind + "\nschemaVersion: " + version + "\n"

		path := writeTempYAML(t, content, "yaml")

		_, err := New(path, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_UnsupportedExtension(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate files with non-YAML extensions
		ext := rapid.SampledFrom([]string{
			"json",
			"toml",
			"xml",
			"ini",
			"conf",
		}).Draw(t, "ext")

		content := "kind: App\nschemaVersion: v1\n"
		path := writeTempYAML(t, content, ext)

		_, err := New(path, testVersions())
		assertErrorFormat(t, err)
	})
}

func TestProperty_ErrorFormat_MapstructureDecodeFailure(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate YAML that parses but has fields that can't be decoded
		// into the target struct (e.g., a map where a string is expected)
		content := "kind: App\nschemaVersion: v1\ncount:\n  nested: value\n"

		path := writeTempYAML(t, content, "yaml")

		_, err := New(path, testVersions())
		// mapstructure may or may not error depending on the input;
		// if it does error, verify the format
		assertErrorFormat(t, err)
	})
}
