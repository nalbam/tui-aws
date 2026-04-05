#!/bin/bash
# SessionStart hook: load project context at the beginning of each session.
# Outputs key project state to help Claude understand current status.

cd "$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0

echo "=== tui-aws Session Context ==="

# Current branch and recent commits
BRANCH=$(git branch --show-current 2>/dev/null)
echo "Branch: $BRANCH"

LAST_COMMIT=$(git log --oneline -1 2>/dev/null)
echo "Last commit: $LAST_COMMIT"

# Uncommitted changes summary
CHANGES=$(git status --porcelain 2>/dev/null | wc -l | tr -d ' ')
if [ "$CHANGES" -gt 0 ]; then
    echo "Uncommitted changes: $CHANGES file(s)"
    git status --porcelain 2>/dev/null | head -10
fi

# Go build status
if command -v go &>/dev/null; then
    echo ""
    echo "Go version: $(go version 2>/dev/null | awk '{print $3}')"
    # Quick vet check (non-blocking)
    VET_OUTPUT=$(go vet ./... 2>&1)
    if [ -n "$VET_OUTPUT" ]; then
        echo "go vet issues detected:"
        echo "$VET_OUTPUT" | head -5
    fi
fi

# Tab count from shared/tab.go
TAB_COUNT=$(grep -c 'Tab[A-Z]' internal/ui/shared/tab.go 2>/dev/null || echo "?")
echo ""
echo "Tabs registered: ~$TAB_COUNT"
echo "Source files: $(find internal/ -name '*.go' 2>/dev/null | wc -l | tr -d ' ')"
echo "=== End Context ==="
