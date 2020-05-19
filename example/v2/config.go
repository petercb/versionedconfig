package v1

import (
	"github.com/petercb/versionedconfig"
)

const (
	Version string = "v2"
	Kind    string = "Config"
)

type ExampleConfig struct {
	Kind          string
	SchemaVersion string
	Metadata      Metadata
	Spec          Spec
}

type Metadata struct {
	Name string
}

type Spec struct {
	Foo string
	Bar int
}

func (obj *ExampleConfig) GetKind() string {
	return obj.Kind
}

func (obj *ExampleConfig) GetVersion() string {
	return obj.SchemaVersion
}

func NewConfig() versionedconfig.VersionedConfig {
	return new(ExampleConfig)
}
