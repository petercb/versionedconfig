// Package v2 defines the v2 schema for the example config.
package v2

import (
	"github.com/petercb/versionedconfig"
	v1 "github.com/petercb/versionedconfig/example/v1"
)

// Version and Kind identify this schema version.
const (
	Version string = "v2"
	Kind    string = "Config"
)

// ExampleConfig is the v2 configuration struct.
type ExampleConfig struct {
	Kind          string
	SchemaVersion string
	Metadata      Metadata
	Spec          Spec
}

// Metadata holds optional metadata for the config.
type Metadata struct {
	Name string
}

// Spec holds the v2 config specification fields.
type Spec struct {
	Foo string
	Bar int
}

// GetKind returns the config kind.
func (obj *ExampleConfig) GetKind() string {
	return obj.Kind
}

// GetVersion returns the schema version.
func (obj *ExampleConfig) GetVersion() string {
	return obj.SchemaVersion
}

// NewConfig is the factory function for creating a v2 ExampleConfig.
func NewConfig() versionedconfig.VersionedConfig {
	return new(ExampleConfig)
}

// UpgradeFromV1 converts a v1 ExampleConfig to a v2 ExampleConfig.
func UpgradeFromV1(cfg versionedconfig.VersionedConfig) (versionedconfig.VersionedConfig, error) {
	old, ok := cfg.(*v1.ExampleConfig)
	if !ok {
		return nil, nil
	}

	return &ExampleConfig{
		Kind:          old.Kind,
		SchemaVersion: Version,
		Metadata:      Metadata{Name: ""},
		Spec: Spec{
			Foo: old.Spec.Foo,
			Bar: old.Spec.Bar,
		},
	}, nil
}
