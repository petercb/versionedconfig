# CLAUDE.md

This file provides context for Claude Code and CI-integrated AI tools working on this repository.

## Product Summary

`versionedconfig` is a Go library for loading and parsing versioned configuration files. Inspired by `k8s.io/apimachinery/pkg/runtime.Scheme`, it provides the same kind+version → Go type dispatch pattern as a lightweight, standalone module with no Kubernetes dependencies.

Key capabilities:

- Read config from file path, URL, or stdin
- Detect schema version from `kind` + `schemaVersion` fields
- Dispatch to a registered factory function for the matching version
- Decode raw YAML into typed Go structs via `mapstructure`
- Provide a simple `VersionedConfig` interface that all config structs implement

## Tech Stack

- **Language:** Go 1.23+
- **Module:** `github.com/petercb/versionedconfig`
- **Dependencies:**
  - `github.com/mitchellh/mapstructure` — Decode generic maps into typed structs
  - `gopkg.in/yaml.v2` — YAML parsing
  - `pgregory.net/rapid` — Property-based testing

## Build & Test Commands

```sh
# Build
go build ./...

# Test
go test ./...

# Test with race detection
go test -race ./...

# Test with coverage
go test -cover ./...

# Lint (golangci-lint v2 required)
golangci-lint run ./...

# Tidy modules
go mod tidy
```

No Makefile or task runner — standard `go` toolchain only.

## Project Structure

```
versionedconfig/
├── types.go              # Core interfaces: VersionedConfig, Version, Versions
├── versionedconfig.go    # Public API: New(), unmarshalConfiguration()
├── util.go               # Internal helpers: file reading, URL detection, download
├── *_test.go             # Unit and property-based tests
├── go.mod / go.sum       # Module definition and dependency lock
├── example/              # Runnable usage example
│   ├── main.go           # Demonstrates registering versions and loading config
│   ├── test.yaml         # Sample YAML config file (kind: Config, schemaVersion: v1)
│   ├── v1/config.go      # v1 schema struct + factory
│   └── v2/config.go      # v2 schema struct + factory (adds Metadata)
└── testdata/rapid/       # Property-based test artifacts
```

## Architecture Pattern

1. **Consumer registers versions** — builds a `Versions` slice of `{SchemaVersion, Kind, Factory}`.
2. **Library reads & parses** — `New(filename, versions)` reads the file, unmarshals YAML into a generic map, extracts `kind` and `schemaVersion`.
3. **Factory dispatch** — looks up the matching factory, instantiates the target struct, decodes the map into it via `mapstructure`.
4. **Consumer uses typed config** — receives a concrete struct behind the `VersionedConfig` interface.

## Conventions

- Each schema version lives in its own sub-package (e.g., `v1/`, `v2/`) exporting `Version`, `Kind`, and `NewConfig`.
- Config structs must implement `VersionedConfig` (i.e., `GetKind()` and `GetVersion()`).
- The library currently only supports YAML config files (`.yaml` / `.yml`).

## Coding Style

- **Formatting:** `gofumpt` and `goimports` are mandatory (enforced by golangci-lint)
- **Immutability:** Prefer returning new objects over mutating existing ones
- **Error handling:** Always wrap errors with context; never silently swallow errors
- **Interfaces:** Accept interfaces, return structs; keep interfaces small (1-3 methods)
- **File size:** 200-400 lines typical, 800 max
- **Functions:** <50 lines per function, no deep nesting (>4 levels)

## Testing

- Use table-driven tests
- Always run with `-race` flag
- Property-based tests use `pgregory.net/rapid`
- Target 80%+ coverage

## Linting

The project uses golangci-lint v2 with strict checks including:

- `govet` with shadow detection
- `errcheck` (unchecked error returns)
- `gocognit`, `gocritic`, `revive`, `misspell`, `unconvert`, `unparam`, `unused`

**Do not modify `.golangci.yaml` without explicit permission.**

## Common Pitfalls

- `defer os.RemoveAll(dir)` — must use `defer func() { _ = os.RemoveAll(dir) }()` to satisfy errcheck
- Variable shadowing with `:=` inside `if` blocks — use a different variable name when outer scope already declares `err`
- Always run `golangci-lint run ./...` after code changes before considering work complete

## Git Workflow

Commit message format:

```
<type>: <description>

<optional body>
```

Types: feat, fix, refactor, docs, test, chore, perf, ci
