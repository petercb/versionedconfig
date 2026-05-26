# Tech Stack

## Language & Runtime
- Go (minimum 1.14)
- Module path: `github.com/petercb/versionedconfig`

## Dependencies
| Package | Purpose |
|---------|---------|
| `github.com/mitchellh/mapstructure` | Decode generic maps into typed structs |
| `github.com/pkg/errors` | Error wrapping with context |
| `gopkg.in/yaml.v2` | YAML parsing |

## Build & Test Commands

```sh
# Build the library
go build ./...

# Run tests
go test ./...

# Run the example
go run ./example/

# Tidy modules
go mod tidy
```

## Notes
- No Makefile or task runner; standard `go` toolchain is used directly.
- No test files exist yet in the repository.
- The library currently only supports YAML config files (`.yaml` / `.yml`).
