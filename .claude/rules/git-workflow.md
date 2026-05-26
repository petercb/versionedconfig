# Git Workflow

## Commit Message Format
```
<type>: <description>

<optional body>
```

Types: feat, fix, refactor, docs, test, chore, perf, ci

## Pull Request Workflow

When creating PRs:
1. Analyze full commit history (not just latest commit)
2. Use `git diff [base-branch]...HEAD` to see all changes
3. Draft comprehensive PR summary
4. Include test plan
5. Push with `-u` flag if new branch

## Development Workflow

1. **Plan First** — Identify dependencies and risks, break down into phases
2. **TDD Approach** — Write tests first (RED), implement (GREEN), refactor (IMPROVE)
3. **Code Review** — Address critical and high issues before merge
4. **Commit & Push** — Detailed commit messages, conventional commits format
