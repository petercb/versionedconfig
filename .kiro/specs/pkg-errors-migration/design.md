# Design Document

## Overview

This design covers the migration of `versionedconfig` from `github.com/pkg/errors` to Go standard library error handling. The migration is a mechanical transformation of 8 call sites across two files (`versionedconfig.go` and `util.go`), replacing deprecated patterns with idiomatic Go 1.22 equivalents while preserving all existing behavior.

## Architecture

The migration does not alter the library's architecture. The existing structure remains:

```
versionedconfig.go  →  Public API (New, unmarshalConfiguration)
util.go             →  Internal helpers (readConfiguration, download, getConfigType, isURL)
types.go            →  Interfaces (unchanged, no error handling)
```

The only changes are to import statements and error construction/wrapping calls within `versionedconfig.go` and `util.go`.

## Components and Interfaces

### Affected Files

| File | Call Sites | Changes |
|------|-----------|---------|
| `versionedconfig.go` | 7 | Replace import, convert 2× `errors.New`, 2× `errors.Errorf`, 3× `errors.Wrap` |
| `util.go` | 1 | Replace import, convert 1× `errors.New` |
| `go.mod` | — | Remove `github.com/pkg/errors v0.9.1` from require block |
| `go.sum` | — | Remove pkg/errors entries |

### Call Site Mapping

#### versionedconfig.go

| Line | Before | After |
|------|--------|-------|
| 1 | `import "github.com/pkg/errors"` | `import ("errors"; "fmt")` |
| 2 | `errors.New("missing kind")` | `errors.New("missing kind")` |
| 3 | `errors.New("missing schemaVersion")` | `errors.New("missing schemaVersion")` |
| 4 | `errors.Errorf("unknown schema version %s/%s", kind, schemaVersion)` | `fmt.Errorf("unknown schema version %s/%s", kind, schemaVersion)` |
| 5 | `errors.Wrap(err, "parse config failure")` | `fmt.Errorf("parse config: %w", err)` |
| 6 | `errors.Wrap(err, "read config")` | `fmt.Errorf("read config: %w", err)` |
| 7 | `errors.Wrap(err, "unmarshal config")` | `fmt.Errorf("unmarshal config: %w", err)` |
| 8 | `errors.Errorf("Unsupported config file type: %s", cfType)` | `fmt.Errorf("unsupported config file type: %s", cfType)` |

#### util.go

| Line | Before | After |
|------|--------|-------|
| 1 | `import "github.com/pkg/errors"` | `import "errors"` |
| 2 | `errors.New("filename not specified")` | `errors.New("filename not specified")` |

### Error Message Refinements

Messages are adjusted to follow Go conventions (lowercase, no trailing punctuation, concise context):

| Before | After | Rationale |
|--------|-------|-----------|
| `"parse config failure"` | `"parse config"` | Remove noun "failure"; the error itself signals failure |
| `"Unsupported config file type: %s"` | `"unsupported config file type: %s"` | Lowercase first letter |

All other messages already conform to Go conventions.

### Public API

No interface changes. The public API signature remains:

```go
func New(filename string, schemaVersions Versions) (VersionedConfig, error)
```

Error behavior is preserved:
- Callers receive `error` values from the same code paths
- Wrapped errors support `errors.Is()` and `errors.As()` via the `%w` verb
- `errors.Unwrap()` traverses the chain identically to `pkg/errors.Cause()`

## Data Models

No data model changes. This migration affects only error construction, not data flow.

## Error Handling

### Wrapping Strategy

All `errors.Wrap` calls become `fmt.Errorf("<context>: %w", err)`. This preserves:
- The error chain (unwrappable via `errors.Unwrap`)
- Compatibility with `errors.Is` and `errors.As`
- Human-readable context when printing the error

### New Error Strategy

`errors.New` calls remain as `errors.New` from the stdlib `errors` package. These are sentinel-style errors with no wrapping.

### Format Error Strategy

`errors.Errorf` calls that don't wrap an existing error become plain `fmt.Errorf`. The one case that wraps (`unknown schema version`) does not wrap an existing error — it creates a new formatted error, so no `%w` verb is needed.

### Post-Migration Import Layout

**versionedconfig.go:**
```go
import (
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v2"
)
```

**util.go:**
```go
import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)
```

## Testing Strategy

### Unit Tests (Example-Based)

- Verify specific error messages match expected strings after migration
- Verify all existing tests pass unchanged (behavior preservation)
- Verify `go build ./...` succeeds

### Property Tests

- **Error chain preservation**: Generate various I/O and parse errors, inject them into wrapping code paths, verify `errors.Is` finds the original cause
- **Error message format**: Trigger all error paths with varied inputs, verify each error message starts lowercase with no trailing punctuation
- **Invalid input rejection**: Generate invalid inputs (empty filenames, YAML missing required fields, unknown versions, unsupported extensions), verify non-nil error returned
- **Valid input acceptance**: Generate valid YAML content with registered kind/schemaVersion, verify nil error and non-nil config returned

### Smoke Tests

- `go mod tidy` does not reintroduce `github.com/pkg/errors`
- No import of `github.com/pkg/errors` in any `.go` file
- No new third-party dependencies added to `go.mod`

## Correctness Properties

*A property is a characteristic or behavior that should hold true across all valid executions of a system — essentially, a formal statement about what the system should do. Properties serve as the bridge between human-readable specifications and machine-verifiable correctness guarantees.*

### Property 1: Error chain preservation

*For any* error returned by a wrapping code path (read failure, unmarshal failure, decode failure), calling `errors.Is(returnedErr, originalCause)` SHALL return true, confirming the original error is preserved in the chain.

**Validates: Requirements 1.2, 2.3**

### Property 2: Error message format conventions

*For any* error returned by the library, the error message string SHALL start with a lowercase letter and SHALL NOT end with punctuation (period, exclamation, or question mark).

**Validates: Requirements 2.1**

### Property 3: Invalid input error preservation

*For any* input that triggers an error path (empty filename, missing kind field, missing schemaVersion field, unknown schema version, invalid YAML, unsupported file type), the library SHALL return a non-nil error.

**Validates: Requirements 4.1**

### Property 4: Valid input success preservation

*For any* valid YAML file containing a registered kind and schemaVersion with well-formed fields, the library SHALL return a non-nil VersionedConfig and a nil error.

**Validates: Requirements 4.2**
