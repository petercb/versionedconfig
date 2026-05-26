---
paths:
  - "**/*.go"
---
# Security Guidelines

## Mandatory Security Checks

Before ANY commit:
- [ ] No hardcoded secrets (API keys, passwords, tokens)
- [ ] All user inputs validated
- [ ] Error messages don't leak sensitive data

## Secret Management

- NEVER hardcode secrets in source code
- ALWAYS use environment variables or a secret manager
- Validate that required secrets are present at startup

## Security Scanning

- Use **gosec** for static security analysis:
  ```bash
  gosec ./...
  ```
- Use **govulncheck** for dependency vulnerabilities:
  ```bash
  govulncheck ./...
  ```

## Go-Specific Security

- Validate all URL inputs before fetching (this library downloads configs from URLs)
- Use `filepath.Clean` for file paths
- Set timeouts on HTTP clients
- Avoid `unsafe` package unless absolutely necessary
