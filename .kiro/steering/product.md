# Product Summary

`versionedconfig` is a Go library for loading and parsing versioned configuration files. Inspired by `k8s.io/apimachinery/pkg/runtime.Scheme`, it provides the same kind+version → Go type dispatch pattern as a lightweight, standalone module with no Kubernetes dependencies.

It allows consumers to define multiple schema versions for a config "kind" and automatically routes deserialization to the correct struct based on `kind` and `schemaVersion` fields in the YAML file.

Key capabilities:
- Read config from file path, URL, or stdin
- Detect schema version from `kind` + `schemaVersion` fields
- Dispatch to a registered factory function for the matching version
- Decode raw YAML into typed Go structs via `mapstructure`
- Provide a simple `VersionedConfig` interface that all config structs implement
