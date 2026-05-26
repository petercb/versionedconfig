// Package v1 defines the v1 schema for the example config.
package v1

import (
	"github.com/petercb/versionedconfig"
)

// Version and Kind identify this schema version.
const (
	Version string = "v1"
	Kind    string = "Config"
)

// ExampleConfig is the v1 configuration struct.
type ExampleConfig struct {
	Kind          string
	SchemaVersion string
	Spec          Spec
}

// Spec holds the v1 config specification fields.
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

// NewConfig is the factory function for creating a v1 ExampleConfig.
func NewConfig() versionedconfig.VersionedConfig {
	return new(ExampleConfig)
}
