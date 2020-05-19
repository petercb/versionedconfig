// Package VersionedConfig exposes functionality to define and use
// different versions of config schemas

package versionedconfig

import (
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// New instantiates a new VersionedConfig interface
func New(filename string, schemaVersions Versions) (VersionedConfig, error) {
	cfg, err := unmarshalConfiguration(filename)
	if err != nil {
		return nil, err
	}

	kind, present := cfg["kind"]
	if !present {
		return nil, errors.New("missing kind")
	}

	schemaVersion, present := cfg["schemaVersion"]
	if !present {
		return nil, errors.New("missing schemaVersion")
	}

	factory, present := schemaVersions.Find(kind.(string), schemaVersion.(string))
	if !present {
		return nil, errors.Errorf("unknown schema version %s/%s", kind, schemaVersion)
	}

	result := factory()

	err = mapstructure.Decode(cfg, &result)
	if err != nil {
		return nil, errors.Wrap(err, "parse config failure")
	}

	return result, nil
}

// UnmarshalConfiguration reads a configuration file
// It returns a map of the contents and an error value on failure
func unmarshalConfiguration(filename string) (map[string]interface{}, error) {
	buf, err := readConfiguration(filename)
	if err != nil {
		return nil, errors.Wrap(err, "read config")
	}

	output := make(map[string]interface{})

	// TODO: Maybe support other file formats for config in the future
	// For now, only support YAML
	cfType := getConfigType(filename)
	switch cfType {
	case "yaml":
		if err = yaml.Unmarshal(buf, &output); err != nil {
			return nil, errors.Wrap(err, "unmarshal config")
		}
	default:
		return nil, errors.Errorf("Unsupported config file type: %s", cfType)
	}

	return output, nil
}
