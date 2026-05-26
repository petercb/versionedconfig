package versionedconfig

import "testing"

// stubConfig is a minimal VersionedConfig for testing.
type stubConfig struct {
	kind    string
	version string
}

func (s *stubConfig) GetKind() string    { return s.kind }
func (s *stubConfig) GetVersion() string { return s.version }

func newStubFactory(kind, version string) func() VersionedConfig {
	return func() VersionedConfig {
		return &stubConfig{kind: kind, version: version}
	}
}

func TestVersions_Find(t *testing.T) {
	versions := Versions{
		{SchemaVersion: "v1", Kind: "Config", Factory: newStubFactory("Config", "v1")},
		{SchemaVersion: "v2", Kind: "Config", Factory: newStubFactory("Config", "v2")},
		{SchemaVersion: "v1", Kind: "Secret", Factory: newStubFactory("Secret", "v1")},
	}

	tests := []struct {
		name          string
		kind          string
		schemaVersion string
		wantFound     bool
		wantKind      string
		wantVersion   string
	}{
		{
			name:          "find v1 Config",
			kind:          "Config",
			schemaVersion: "v1",
			wantFound:     true,
			wantKind:      "Config",
			wantVersion:   "v1",
		},
		{
			name:          "find v2 Config",
			kind:          "Config",
			schemaVersion: "v2",
			wantFound:     true,
			wantKind:      "Config",
			wantVersion:   "v2",
		},
		{
			name:          "find v1 Secret",
			kind:          "Secret",
			schemaVersion: "v1",
			wantFound:     true,
			wantKind:      "Secret",
			wantVersion:   "v1",
		},
		{
			name:          "unknown kind",
			kind:          "Unknown",
			schemaVersion: "v1",
			wantFound:     false,
		},
		{
			name:          "unknown version",
			kind:          "Config",
			schemaVersion: "v99",
			wantFound:     false,
		},
		{
			name:          "empty kind",
			kind:          "",
			schemaVersion: "v1",
			wantFound:     false,
		},
		{
			name:          "empty version",
			kind:          "Config",
			schemaVersion: "",
			wantFound:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factory, found := versions.Find(tt.kind, tt.schemaVersion)
			if found != tt.wantFound {
				t.Fatalf("Find(%q, %q) found = %v, want %v", tt.kind, tt.schemaVersion, found, tt.wantFound)
			}
			if !found {
				if factory != nil {
					t.Errorf("Find(%q, %q) returned non-nil factory when not found", tt.kind, tt.schemaVersion)
				}
				return
			}
			cfg := factory()
			if cfg.GetKind() != tt.wantKind {
				t.Errorf("factory().GetKind() = %q, want %q", cfg.GetKind(), tt.wantKind)
			}
			if cfg.GetVersion() != tt.wantVersion {
				t.Errorf("factory().GetVersion() = %q, want %q", cfg.GetVersion(), tt.wantVersion)
			}
		})
	}
}

func TestVersions_Find_EmptySlice(t *testing.T) {
	versions := Versions{}
	factory, found := versions.Find("Config", "v1")
	if found {
		t.Error("Find on empty Versions should return false")
	}
	if factory != nil {
		t.Error("Find on empty Versions should return nil factory")
	}
}
