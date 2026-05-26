---
paths:
  - "**/*.go"
  - "**/go.mod"
  - "**/go.sum"
---
# Go Coding Style

## Formatting

- **gofumpt** and **goimports** are mandatory — no style debates
- Enforced by golangci-lint; do not modify `.golangci.yaml` without explicit permission

## Design Principles

- Accept interfaces, return structs
- Keep interfaces small (1-3 methods)
- Define interfaces where they are used, not where they are implemented

## Error Handling

- Always wrap errors with context
- Never silently swallow errors
- Handle errors explicitly at every level
- Provide clear error messages

## Immutability

ALWAYS create new objects, NEVER mutate existing ones. Rationale: prevents hidden side effects, makes debugging easier, enables safe concurrency.

## File Organization

- 200-400 lines typical, 800 max
- Functions <50 lines, no deep nesting (>4 levels)
- High cohesion, low coupling
- Organize by feature/domain

## Input Validation

- Validate all inputs at system boundaries
- Fail fast with clear error messages
- Never trust external data
