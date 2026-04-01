---
name: code-review
description: Review changed code with confidence-based scoring. Use when asked to review code, check for bugs, or audit changes.
---

# Code Review Skill

Review changed code with confidence-based scoring to filter false positives.

## Review Scope
By default, review unstaged changes from `git diff`. The user may specify different files or scope.

## Review Criteria
- Import patterns and module boundaries
- Framework conventions (Bubble Tea, aws-sdk-go-v2)
- Bug detection: logic errors, nil handling, race conditions
- Security vulnerabilities (OWASP Top 10)
- Code duplication and unnecessary complexity

## Confidence Scoring
Rate each issue 0-100. Only report issues with confidence >= 75.
- 75-89: Verified real issue, report with fix suggestion
- 90-100: Confirmed critical issue, must report

## Output Format
```
### [CRITICAL|IMPORTANT] <issue title> (confidence: XX)
**File:** `path/to/file.ext:line`
**Issue:** Clear description
**Fix:** Concrete code suggestion
```
