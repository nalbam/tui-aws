# Code Review

Review the current diff for bugs, style issues, and project convention violations.

## Steps

1. Run `git diff` to see all unstaged changes
2. Run `git diff --cached` to see staged changes
3. Read CLAUDE.md for project conventions
4. For each changed file:
   - Check Go conventions: error handling, naming, imports
   - Check Bubble Tea patterns: Model-View-Update correctness
   - Check SharedState usage: no direct mutation outside Update()
   - Check tab architecture: TabModel interface compliance
   - Look for: race conditions, nil pointer dereferences, resource leaks
   - Verify lipgloss.Width() usage for Unicode-safe rendering
5. Rate each issue with confidence score (0-100)
6. Report only issues with confidence >= 75

## Output Format

For each issue:
```
### [CRITICAL|IMPORTANT] <issue title> (confidence: XX)
**File:** `path/to/file.go:line`
**Issue:** Description
**Fix:** Concrete suggestion
```
