package versionedconfig

import (
	"fmt"
)

// UpgradeFunc is the function signature for upgrading a config from one version
// to the next. It takes a VersionedConfig at the current version and returns a
// new VersionedConfig at the next version in the chain.
type UpgradeFunc func(VersionedConfig) (VersionedConfig, error)

// Upgrader defines the interface for types that can upgrade a VersionedConfig
// through a chain of versions to reach the latest.
type Upgrader interface {
	Upgrade(VersionedConfig) (VersionedConfig, error)
}

// Upgrade takes a VersionedConfig and upgrades it through the registered chain
// to the latest version for its kind. If the config is already at the latest
// version, it is returned as-is. Returns an error if any upgrade step fails,
// if a required UpgradeTo function is missing, or if the kind/version is unknown.
func (obj *Versions) Upgrade(cfg VersionedConfig) (VersionedConfig, error) {
	kind := cfg.GetKind()
	version := cfg.GetVersion()

	// Collect all versions for this kind in order.
	var chain []Version
	for _, v := range *obj {
		if v.Kind == kind {
			chain = append(chain, v)
		}
	}

	if len(chain) == 0 {
		return nil, fmt.Errorf("upgrade: unknown kind %q", kind)
	}

	// Find the starting position in the chain.
	startIdx := -1
	for i, v := range chain {
		if v.SchemaVersion == version {
			startIdx = i
			break
		}
	}

	if startIdx < 0 {
		return nil, fmt.Errorf("upgrade: unknown version %q for kind %q", version, kind)
	}

	// If already at the latest version, return as-is.
	if startIdx == len(chain)-1 {
		return cfg, nil
	}

	// Run upgrade chain from startIdx to the last version.
	current := cfg
	for i := startIdx; i < len(chain)-1; i++ {
		if chain[i].UpgradeTo == nil {
			return nil, fmt.Errorf(
				"upgrade %s/%s: no upgrade function registered",
				kind, chain[i].SchemaVersion,
			)
		}

		next, err := chain[i].UpgradeTo(current)
		if err != nil {
			return nil, fmt.Errorf(
				"upgrade %s to %s: %w",
				chain[i].SchemaVersion, chain[i+1].SchemaVersion, err,
			)
		}

		current = next
	}

	return current, nil
}

// NewWithUpgrade loads a config file and automatically upgrades it to the latest
// registered version for its kind. It combines the behavior of New() with
// sequential version upgrading. If the config is already at the latest version,
// no upgrade is performed.
func NewWithUpgrade(filename string, schemaVersions Versions) (VersionedConfig, error) {
	cfg, err := New(filename, schemaVersions)
	if err != nil {
		return nil, err
	}

	return schemaVersions.Upgrade(cfg)
}
