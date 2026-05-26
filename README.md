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

## Install

```sh
go get github.com/petercb/versionedconfig
```
