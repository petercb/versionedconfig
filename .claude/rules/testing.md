---
paths:
  - "**/*.go"
  - "**/*_test.go"
---
# Testing Requirements

## Framework

Use the standard `go test` with **table-driven tests**.

## Race Detection

Always run with the `-race` flag:
```bash
go test -race ./...
```

## Coverage

Target 80%+ coverage:
```bash
go test -cover ./...
```

## Property-Based Testing

This project uses `pgregory.net/rapid` for property-based tests. Use it for testing invariants across random inputs.

## Test-Driven Development

MANDATORY workflow:
1. Write test first (RED)
2. Run test — it should FAIL
3. Write minimal implementation (GREEN)
4. Run test — it should PASS
5. Refactor (IMPROVE)
6. Verify coverage (80%+)

## Troubleshooting

1. Check test isolation
2. Verify mocks are correct
3. Fix implementation, not tests (unless tests are wrong)
