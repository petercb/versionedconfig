package versionedconfig

import (
	"fmt"
	"strings"
	"testing"
	"unicode"

	"pgregory.net/rapid"
)

// TestProperty_UpgradeChain_ReachesLatest verifies that for any valid upgrade
// chain of length N, upgrading from version 1 always produces a config at
// version N.
func TestProperty_UpgradeChain_ReachesLatest(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate a chain length between 2 and 6.
		chainLen := rapid.IntRange(2, 6).Draw(t, "chainLen")

		// Build a version chain with upgrade functions.
		versions := make(Versions, chainLen)
		for i := range chainLen {
			ver := fmt.Sprintf("v%d", i+1)
			versions[i] = Version{
				SchemaVersion: ver,
				Kind:          "Test",
				Factory: func() VersionedConfig {
					return &upgradeTestConfig{kind: "Test", version: ver}
				},
			}
			if i < chainLen-1 {
				nextVer := fmt.Sprintf("v%d", i+2)
				versions[i].UpgradeTo = func(cfg VersionedConfig) (VersionedConfig, error) {
					return &upgradeTestConfig{
						kind:    cfg.GetKind(),
						version: nextVer,
					}, nil
				}
			}
		}

		input := &upgradeTestConfig{kind: "Test", version: "v1"}
		result, err := versions.Upgrade(input)
		if err != nil {
			t.Fatalf("unexpected upgrade error: %v", err)
		}

		expectedVersion := fmt.Sprintf("v%d", chainLen)
		if result.GetVersion() != expectedVersion {
			t.Fatalf("GetVersion() = %q, want %q", result.GetVersion(), expectedVersion)
		}
	})
}

// TestProperty_UpgradeError_Format verifies that upgrade error messages follow
// conventions: lowercase first character and no trailing punctuation.
func TestProperty_UpgradeError_Format(t *testing.T) {
	rapid.Check(t, func(t *rapid.T) {
		// Generate an error message for the upgrade function to return.
		innerMsg := rapid.StringMatching(`[a-z][a-z0-9 ]{2,20}`).Draw(t, "innerMsg")

		versions := Versions{
			{
				SchemaVersion: "v1", Kind: "Svc",
				Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "Svc", version: "v1"} },
				UpgradeTo: func(_ VersionedConfig) (VersionedConfig, error) {
					return nil, fmt.Errorf("%s", innerMsg)
				},
			},
			{
				SchemaVersion: "v2", Kind: "Svc",
				Factory: func() VersionedConfig { return &upgradeTestConfig{kind: "Svc", version: "v2"} },
			},
		}

		input := &upgradeTestConfig{kind: "Svc", version: "v1"}
		_, err := versions.Upgrade(input)
		if err == nil {
			t.Fatal("expected error")
		}

		msg := err.Error()

		// Error should start with lowercase.
		if len(msg) > 0 && unicode.IsUpper(rune(msg[0])) {
			t.Errorf("error starts with uppercase: %q", msg)
		}

		// Error should not end with punctuation.
		if len(msg) > 0 {
			last := msg[len(msg)-1]
			if strings.ContainsRune(".!?", rune(last)) {
				t.Errorf("error ends with punctuation: %q", msg)
			}
		}
	})
}
