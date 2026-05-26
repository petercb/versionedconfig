# Requirements Document

## Introduction

Migrate the `versionedconfig` library from the deprecated `github.com/pkg/errors` package to Go standard library error handling (`errors` and `fmt` packages). The migration covers 8 call sites across two source files, refines error messages to follow Go conventions (lowercase, no trailing punctuation), and removes the third-party dependency entirely.

## Glossary

- **Library**: The `versionedconfig` Go module (`github.com/petercb/versionedconfig`)
- **pkg/errors**: The third-party package `github.com/pkg/errors` v0.9.1, currently in maintenance mode
- **stdlib errors**: The Go standard library `errors` package
- **fmt.Errorf**: The `fmt.Errorf` function used with the `%w` verb for error wrapping
- **Call Site**: A location in source code where a `pkg/errors` function is invoked
- **Go Error Conventions**: Lowercase error strings, no trailing punctuation, concise context at each wrap point

## Requirements

### Requirement 1

**User Story:** As a library maintainer, I want to replace all `pkg/errors` usage with stdlib equivalents, so that the library has no dependency on a deprecated package.

#### Acceptance Criteria

1. WHEN a call site uses `errors.New` from `pkg/errors`, THE Library SHALL replace the call with `errors.New` from the stdlib `errors` package
2. WHEN a call site uses `errors.Wrap` from `pkg/errors`, THE Library SHALL replace the call with `fmt.Errorf` using the `%w` verb to wrap the original error
3. WHEN a call site uses `errors.Errorf` from `pkg/errors`, THE Library SHALL replace the call with `fmt.Errorf` using the `%w` verb where an error is being wrapped, or plain `fmt.Errorf` otherwise
4. THE Library SHALL contain zero import statements referencing `github.com/pkg/errors` after migration

### Requirement 2

**User Story:** As a library maintainer, I want error messages to follow Go conventions, so that the codebase is idiomatic and consistent.

#### Acceptance Criteria

1. THE Library SHALL use lowercase error message strings with no trailing punctuation at every call site
2. WHEN wrapping an error with `fmt.Errorf`, THE Library SHALL provide concise contextual information describing the operation that failed
3. THE Library SHALL preserve the original error in the wrapping chain so that callers can use `errors.Is` and `errors.As` for inspection

### Requirement 3

**User Story:** As a library maintainer, I want the `github.com/pkg/errors` dependency removed from `go.mod`, so that the module dependency graph is minimal and free of deprecated packages.

#### Acceptance Criteria

1. THE Library SHALL remove `github.com/pkg/errors` from the `require` block in `go.mod`
2. THE Library SHALL remove the corresponding entry from `go.sum`
3. THE Library SHALL pass `go mod tidy` without reintroducing the `github.com/pkg/errors` dependency

### Requirement 4

**User Story:** As a library maintainer, I want the migration to preserve existing behavior, so that consumers of the library are unaffected.

#### Acceptance Criteria

1. THE Library SHALL return non-nil errors from the same code paths that returned non-nil errors before migration
2. THE Library SHALL return nil errors from the same code paths that returned nil errors before migration
3. THE Library SHALL pass all existing tests after migration
4. THE Library SHALL build successfully with `go build ./...` after migration

### Requirement 5

**User Story:** As a library maintainer, I want to use only Go 1.22 stdlib features for error handling, so that no new third-party dependencies are introduced.

#### Acceptance Criteria

1. THE Library SHALL not introduce any new third-party dependencies for error handling
2. THE Library SHALL use only the `errors` and `fmt` packages from the Go standard library for error creation and wrapping
3. WHEN a use case benefits from `errors.Join`, THE Library SHALL use it rather than introducing a third-party alternative
