# Project Structure

```
versionedconfig/
├── types.go              # Core interfaces: VersionedConfig, Version, Versions
├── versionedconfig.go    # Public API: New(), unmarshalConfiguration()
├── util.go               # Internal helpers: file reading, URL detection, download
├── go.mod / go.sum       # Module definition and dependency lock
├── example/              # Runnable usage example
│   ├── main.go           # Demonstrates registering versions and loading config
│   ├── test.yaml         # Sample YAML config file (kind: Config, schemaVersion: v1)
│   ├── v1/config.go      # v1 schema struct + factory
│   └── v2/config.go      # v2 schema struct + factory (adds Metadata)
└── .kiro/steering/       # AI assistant steering rules
```

## Architecture Pattern

1. **Consumer registers versions** — builds a `Versions` slice of `{SchemaVersion, Kind, Factory}`.
2. **Library reads & parses** — `New(filename, versions)` reads the file, unmarshals YAML into a generic map, extracts `kind` and `schemaVersion`.
3. **Factory dispatch** — looks up the matching factory, instantiates the target struct, decodes the map into it via `mapstructure`.
4. **Consumer uses typed config** — receives a concrete struct behind the `VersionedConfig` interface.

## Conventions
- Each schema version lives in its own sub-package (e.g., `v1/`, `v2/`) exporting `Version`, `Kind`, and `NewConfig`.
- Config structs must implement `VersionedConfig` (i.e., `GetKind()` and `GetVersion()`).
- Package-level doc comment is in `versionedconfig.go`.
