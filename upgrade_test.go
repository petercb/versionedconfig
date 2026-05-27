package versionedconfig

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// upgradeTestConfig implements VersionedConfig for upgrade tests.
type upgradeTestConfig struct {
	kind    string
	version string
	data    string
}

func (c *upgradeTestConfig) GetKind() string    { return c.kind }
func (c *upgradeTestConfig) GetVersion() string { return c.version }

func TestVersions_Upgrade(t *testing.T) {
	errUpgrade := errors.New("conversion failed")

	tests := []struct {
		name     string
		versions Versions
		input    VersionedConfig
		wantVer  string
		wantErr  bool
		errMsg   string
	}{
		{
			name: "single step upgrade v1 to v2",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
					UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
						return &upgradeTestConfig{kind: cfg.GetKind(), version: "v2", data: "upgraded"}, nil
					},
				},
				{
					SchemaVersion: "v2", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v1"},
			wantVer: "v2",
		},
		{
			name: "multi step upgrade v1 to v3",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
					UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
						return &upgradeTestConfig{kind: cfg.GetKind(), version: "v2", data: "from-v1"}, nil
					},
				},
				{
					SchemaVersion: "v2", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
					UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
						return &upgradeTestConfig{kind: cfg.GetKind(), version: "v3", data: "from-v2"}, nil
					},
				},
				{
					SchemaVersion: "v3", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v3"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v1"},
			wantVer: "v3",
		},
		{
			name: "already at latest version",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
					UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
						return &upgradeTestConfig{kind: cfg.GetKind(), version: "v2"}, nil
					},
				},
				{
					SchemaVersion: "v2", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v2"},
			wantVer: "v2",
		},
		{
			name: "nil UpgradeTo for non-latest version",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory:   func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
					UpgradeTo: nil,
				},
				{
					SchemaVersion: "v2", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v1"},
			wantErr: true,
			errMsg:  "no upgrade function registered",
		},
		{
			name: "upgrade function returns error",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
					UpgradeTo: func(_ VersionedConfig) (VersionedConfig, error) {
						return nil, errUpgrade
					},
				},
				{
					SchemaVersion: "v2", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v1"},
			wantErr: true,
			errMsg:  "upgrade v1 to v2",
		},
		{
			name: "unknown kind",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
				},
			},
			input:   &upgradeTestConfig{kind: "Unknown", version: "v1"},
			wantErr: true,
			errMsg:  "unknown kind",
		},
		{
			name: "unknown version",
			versions: Versions{
				{
					SchemaVersion: "v1", Kind: "App",
					Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
				},
			},
			input:   &upgradeTestConfig{kind: "App", version: "v99"},
			wantErr: true,
			errMsg:  "unknown version",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.versions.Upgrade(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if tt.errMsg != "" && !containsStr(err.Error(), tt.errMsg) {
					t.Errorf("error = %q, want substring %q", err.Error(), tt.errMsg)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result.GetVersion() != tt.wantVer {
				t.Errorf("GetVersion() = %q, want %q", result.GetVersion(), tt.wantVer)
			}
		})
	}
}

func TestVersions_Upgrade_ErrorChainPreserved(t *testing.T) {
	underlying := errors.New("disk full")
	versions := Versions{
		{
			SchemaVersion: "v1", Kind: "App",
			Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
			UpgradeTo: func(_ VersionedConfig) (VersionedConfig, error) {
				return nil, underlying
			},
		},
		{
			SchemaVersion: "v2", Kind: "App",
			Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
		},
	}

	_, err := versions.Upgrade(&upgradeTestConfig{kind: "App", version: "v1"})
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, underlying) {
		t.Errorf("errors.Is(err, underlying) = false; err = %v", err)
	}
}

func TestNewWithUpgrade_Success(t *testing.T) {
	yaml := "kind: App\nschemaVersion: v1\nname: myapp\ncount: 10\n"
	path := writeUpgradeYAML(t, yaml)

	versions := Versions{
		{
			SchemaVersion: "v1", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
			UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
				tc := cfg.(*testConfig)
				return &testConfig{
					Kind:          tc.Kind,
					SchemaVersion: "v2",
					Name:          tc.Name + "-upgraded",
					Count:         tc.Count * 2,
				}, nil
			},
		},
		{
			SchemaVersion: "v2", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
		},
	}

	result, err := NewWithUpgrade(path, versions)
	if err != nil {
		t.Fatalf("NewWithUpgrade() error: %v", err)
	}

	tc, ok := result.(*testConfig)
	if !ok {
		t.Fatalf("expected *testConfig, got %T", result)
	}
	if tc.SchemaVersion != "v2" {
		t.Errorf("SchemaVersion = %q, want %q", tc.SchemaVersion, "v2")
	}
	if tc.Name != "myapp-upgraded" {
		t.Errorf("Name = %q, want %q", tc.Name, "myapp-upgraded")
	}
	if tc.Count != 20 {
		t.Errorf("Count = %d, want %d", tc.Count, 20)
	}
}

func TestNewWithUpgrade_AlreadyLatest(t *testing.T) {
	yaml := "kind: App\nschemaVersion: v2\nname: latest\ncount: 5\n"
	path := writeUpgradeYAML(t, yaml)

	versions := Versions{
		{
			SchemaVersion: "v1", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
			UpgradeTo: func(_ VersionedConfig) (VersionedConfig, error) {
				return nil, errors.New("should not be called")
			},
		},
		{
			SchemaVersion: "v2", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
		},
	}

	result, err := NewWithUpgrade(path, versions)
	if err != nil {
		t.Fatalf("NewWithUpgrade() error: %v", err)
	}

	tc := result.(*testConfig)
	if tc.SchemaVersion != "v2" {
		t.Errorf("SchemaVersion = %q, want %q", tc.SchemaVersion, "v2")
	}
	if tc.Name != "latest" {
		t.Errorf("Name = %q, want %q", tc.Name, "latest")
	}
}

func TestNewWithUpgrade_FileError(t *testing.T) {
	versions := Versions{
		{SchemaVersion: "v1", Kind: "App", Factory: func() VersionedConfig { return &testConfig{} }},
	}

	_, err := NewWithUpgrade("/nonexistent/config.yaml", versions)
	if err == nil {
		t.Fatal("expected error for nonexistent file")
	}
}

func TestVersions_Upgrade_NilReturn(t *testing.T) {
	versions := Versions{
		{
			SchemaVersion: "v1", Kind: "App",
			Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v1"} },
			UpgradeTo: func(_ VersionedConfig) (VersionedConfig, error) {
				return nil, nil
			},
		},
		{
			SchemaVersion: "v2", Kind: "App",
			Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "App", version: "v2"} },
		},
	}

	_, err := versions.Upgrade(&upgradeTestConfig{kind: "App", version: "v1"})
	if err == nil {
		t.Fatal("expected error when upgrade function returns nil")
	}
	if !containsStr(err.Error(), "upgrade function returned nil") {
		t.Errorf("error = %q, want substring %q", err.Error(), "upgrade function returned nil")
	}
}

func TestNewWithUpgrade_MultiStep(t *testing.T) {
	yaml := "kind: App\nschemaVersion: v1\nname: myapp\ncount: 5\n"
	path := writeUpgradeYAML(t, yaml)

	versions := Versions{
		{
			SchemaVersion: "v1", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
			UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
				tc := cfg.(*testConfig)
				return &testConfig{
					Kind:          tc.Kind,
					SchemaVersion: "v2",
					Name:          tc.Name + "-v2",
					Count:         tc.Count + 1,
				}, nil
			},
		},
		{
			SchemaVersion: "v2", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
			UpgradeTo: func(cfg VersionedConfig) (VersionedConfig, error) {
				tc := cfg.(*testConfig)
				return &testConfig{
					Kind:          tc.Kind,
					SchemaVersion: "v3",
					Name:          tc.Name + "-v3",
					Count:         tc.Count + 1,
				}, nil
			},
		},
		{
			SchemaVersion: "v3", Kind: "App",
			Factory: func() VersionedConfig { return &testConfig{} },
		},
	}

	result, err := NewWithUpgrade(path, versions)
	if err != nil {
		t.Fatalf("NewWithUpgrade() error: %v", err)
	}

	tc, ok := result.(*testConfig)
	if !ok {
		t.Fatalf("expected *testConfig, got %T", result)
	}
	if tc.SchemaVersion != "v3" {
		t.Errorf("SchemaVersion = %q, want %q", tc.SchemaVersion, "v3")
	}
	if tc.Name != "myapp-v2-v3" {
		t.Errorf("Name = %q, want %q", tc.Name, "myapp-v2-v3")
	}
	if tc.Count != 7 {
		t.Errorf("Count = %d, want %d", tc.Count, 7)
	}
}

func writeUpgradeYAML(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write test yaml: %v", err)
	}
	return path
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && contains(s, substr))
}

func contains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
