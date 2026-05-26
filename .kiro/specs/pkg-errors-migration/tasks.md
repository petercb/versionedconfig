# Implementation Plan: pkg-errors-migration

## Overview

Migrate all 8 `pkg/errors` call sites to Go stdlib equivalents (`errors` and `fmt`), refine error messages to follow Go conventions, remove the deprecated dependency, and validate correctness with property-based tests.

## Tasks

- [x] 1. Migrate versionedconfig.go error handling
  - [x] 1.1 Replace imports and convert error calls in versionedconfig.go
    - Remove `"github.com/pkg/errors"` from the import block
    - Add `"errors"` and `"fmt"` to the import block (grouped with stdlib imports)
    - Convert `errors.New("missing kind")` → `errors.New("missing kind")` (stdlib)
    - Convert `errors.New("missing schemaVersion")` → `errors.New("missing schemaVersion")` (stdlib)
    - Convert `errors.Errorf("unknown schema version %s/%s", kind, schemaVersion)` → `fmt.Errorf("unknown schema version %s/%s", kind, schemaVersion)`
    - Convert `errors.Wrap(err, "parse config failure")` → `fmt.Errorf("parse config: %w", err)`
    - Convert `errors.Wrap(err, "read config")` → `fmt.Errorf("read config: %w", err)`
    - Convert `errors.Wrap(err, "unmarshal config")` → `fmt.Errorf("unmarshal config: %w", err)`
    - Convert `errors.Errorf("Unsupported config file type: %s", cfType)` → `fmt.Errorf("unsupported config file type: %s", cfType)`
    - _Requirements: 1.1, 1.2, 1.3, 1.4, 2.1, 2.2_

- [x] 2. Migrate util.go error handling
  - [x] 2.1 Replace import and convert error call in util.go
    - Remove `"github.com/pkg/errors"` from the import block
    - Add `"errors"` to the import block (grouped with stdlib imports)
    - Convert `errors.New("filename not specified")` → `errors.New("filename not specified")` (stdlib)
    - Verify final import block matches design: `errors`, `io`, `net/http`, `os`, `path/filepath`, `strings`
    - _Requirements: 1.1, 1.4, 2.1_

- [x] 3. Remove pkg/errors dependency from module
  - [x] 3.1 Run go mod tidy to remove pkg/errors
    - Run `go mod tidy` to remove `github.com/pkg/errors` from `go.mod` and `go.sum`
    - Verify `go.mod` no longer contains `github.com/pkg/errors` in the require block
    - Verify `go.sum` no longer contains `github.com/pkg/errors` entries
    - _Requirements: 3.1, 3.2, 3.3_

- [x] 4. Checkpoint - Verify build and existing tests
  - Run `go build ./...` and confirm zero errors
  - Run `go test ./...` and confirm all existing tests pass
  - Ensure all tests pass, ask the user if questions arise.
  - _Requirements: 4.3, 4.4, 5.1, 5.2_

- [x] 5. Write property-based tests
  - [x] 5.1 Write property test for error chain preservation
    - **Property 1: Error chain preservation**
    - Create `versionedconfig_prop_test.go` with rapid-based property tests
    - Generate various I/O and parse errors, inject via temp files with invalid content
    - Verify `errors.Is(returnedErr, originalCause)` returns true for all wrapping paths
    - Use `pgregory.net/rapid` for property-based test generation
    - **Validates: Requirements 1.2, 2.3**

  - [x] 5.2 Write property test for error message format conventions
    - **Property 2: Error message format conventions**
    - Trigger all error paths with varied inputs (empty filenames, bad YAML, unknown versions, unsupported extensions)
    - Verify each error message starts with a lowercase letter
    - Verify no error message ends with `.`, `!`, or `?`
    - **Validates: Requirements 2.1**

  - [x] 5.3 Write property test for invalid input error preservation
    - **Property 3: Invalid input error preservation**
    - Generate invalid inputs: empty filenames, YAML missing `kind`, YAML missing `schemaVersion`, unknown kind/version combos, non-YAML extensions
    - Verify the library returns a non-nil error for every invalid input
    - **Validates: Requirements 4.1**

  - [x] 5.4 Write property test for valid input success preservation
    - **Property 4: Valid input success preservation**
    - Generate valid YAML files with registered kind and schemaVersion
    - Write generated YAML to temp files, call `New()` with matching `Versions`
    - Verify non-nil `VersionedConfig` and nil error returned
    - **Validates: Requirements 4.2**

- [x] 6. Final checkpoint - Verify all tests pass
  - Run `go test ./...` and confirm all tests pass (including property tests)
  - Run `go build ./...` and confirm clean build
  - Verify no import of `github.com/pkg/errors` exists in any `.go` file
  - Ensure all tests pass, ask the user if questions arise.

## Notes

- Tasks marked with `*` are optional and can be skipped for faster MVP
- Each task references specific requirements for traceability
- Checkpoints ensure incremental validation
- Property tests use `pgregory.net/rapid` (the standard Go property testing library)
- The migration is mechanical — no logic changes, only error construction patterns change
- The existing test files (`versionedconfig_test.go`, `util_test.go`, `types_test.go`) serve as regression guards

## Task Dependency Graph

```json
{
  "waves": [
    { "id": 0, "tasks": ["1.1", "2.1"] },
    { "id": 1, "tasks": ["3.1"] },
    { "id": 2, "tasks": ["5.1", "5.2", "5.3", "5.4"] }
  ]
}
```
