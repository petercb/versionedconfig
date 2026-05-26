// Package main demonstrates usage of the versionedconfig library.
package main

import (
	"fmt"

	"github.com/petercb/versionedconfig"
	v1 "github.com/petercb/versionedconfig/example/v1"
	v2 "github.com/petercb/versionedconfig/example/v2"
)

func main() {
	schemaVersions := versionedconfig.Versions{
		{SchemaVersion: v1.Version, Kind: v1.Kind, Factory: v1.NewConfig},
		{SchemaVersion: v2.Version, Kind: v2.Kind, Factory: v2.NewConfig},
	}
	myConfig, err := versionedconfig.New("test.yaml", schemaVersions)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", myConfig)
}
