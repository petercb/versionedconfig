# versionedconfig

Lightweight versioned config file handling for Go, inspired by [`k8s.io/apimachinery/pkg/runtime.Scheme`](https://pkg.go.dev/k8s.io/apimachinery/pkg/runtime#Scheme).

Provides the same kind+version → Go type dispatch pattern without pulling in the Kubernetes dependency tree.

## Features

- Read configuration from a file path, URL, or stdin
- Route deserialization based on `kind` and `schemaVersion` fields
- Register multiple schema versions with factory functions
- Decode YAML into typed Go structs via [`mapstructure`](https://github.com/mitchellh/mapstructure)

## Usage

```go
package main

import (
 "fmt"

 "github.com/petercb/versionedconfig"
 v1 "github.com/petercb/versionedconfig/example/v1"
 v2 "github.com/petercb/versionedconfig/example/v2"
)

func main() {
 schemaVersions := versionedconfig.Versions{
  {v1.Version, v1.Kind, v1.NewConfig},
  {v2.Version, v2.Kind, v2.NewConfig},
 }

 cfg, err := versionedconfig.New("config.yaml", schemaVersions)
 if err != nil {
  panic(err)
 }

 fmt.Printf("%#v\n", cfg)
}
```

Config files must include `kind` and `schemaVersion` at the top level:

```yaml
kind: Config
schemaVersion: v1
spec:
  foo: bar
  bar: 42
```

## Implementing a Schema Version

Each version lives in its own package and exports a factory function:

```go
package v1

import "github.com/petercb/versionedconfig"

const (
 Version string = "v1"
 Kind    string = "Config"
)

type MyConfig struct {
 Kind          string
 SchemaVersion string
 Spec          Spec
}

type Spec struct {
 Foo string
 Bar int
}

func (c *MyConfig) GetKind() string    { return c.Kind }
func (c *MyConfig) GetVersion() string { return c.SchemaVersion }

func NewConfig() versionedconfig.VersionedConfig {
 return new(MyConfig)
}
```

## Automatic Version Upgrades

Register upgrade functions to automatically convert older config files to the
latest version. Each version in the chain specifies an `UpgradeTo` function that
converts it to the next version.

### Defining Upgrade Functions

```go
package v2

import (
 "github.com/petercb/versionedconfig"
 v1 "example/v1"
)

func UpgradeFromV1(cfg versionedconfig.VersionedConfig) (versionedconfig.VersionedConfig, error) {
 old := cfg.(*v1.MyConfig)
 return &MyConfigV2{
  Kind:          old.Kind,
  SchemaVersion: "v2",
  Spec:          Spec{Foo: old.Spec.Foo, Bar: old.Spec.Bar},
 }, nil
}
```

### Registering Upgrades

Add the `UpgradeTo` field when registering versions. The versions slice must be
ordered from oldest to newest for each kind:

```go
schemaVersions := versionedconfig.Versions{
 {SchemaVersion: v1.Version, Kind: v1.Kind, Factory: v1.NewConfig, UpgradeTo: v2.UpgradeFromV1},
 {SchemaVersion: v2.Version, Kind: v2.Kind, Factory: v2.NewConfig, UpgradeTo: v3.UpgradeFromV2},
 {SchemaVersion: v3.Version, Kind: v3.Kind, Factory: v3.NewConfig}, // latest, no UpgradeTo
}
```

### Loading with Automatic Upgrade

Use `NewWithUpgrade` to load a config file and automatically upgrade it to the
latest registered version:

```go
cfg, err := versionedconfig.NewWithUpgrade("config.yaml", schemaVersions)
// cfg is now at the latest version, even if the file was v1
```

If the file is already at the latest version, it is returned without modification.
If any upgrade step fails, the error is wrapped with context identifying which
version transition failed.

### Upgrading an Existing Config

You can also upgrade a config that was already loaded with `New()`:

```go
cfg, _ := versionedconfig.New("config.yaml", schemaVersions)
upgraded, err := schemaVersions.Upgrade(cfg)
```

## Install

```sh
go get github.com/petercb/versionedconfig
```
