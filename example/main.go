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
	myConfig, err := versionedconfig.New("test.yaml", schemaVersions)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v", myConfig)
}
