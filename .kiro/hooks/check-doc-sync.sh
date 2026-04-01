#!/bin/bash
# Detect documentation sync needs after file changes.
# Triggered by PostToolUse (fs_write) events.

FILE_PATH="${1:-}"
[ -z "$FILE_PATH" ] && exit 0

# Alert if steering rules might need update
if [[ "$FILE_PATH" == internal/aws/* ]]; then
    echo "[doc-sync] AWS module changed. Check .kiro/steering/aws-sdk.md if new service added."
fi

if [[ "$FILE_PATH" == internal/ui/tab_* ]]; then
    echo "[doc-sync] Tab changed. Check .kiro/steering/tab-architecture.md if new tab added."
fi

# Alert if no ADRs exist when architecture files change
if [[ "$FILE_PATH" == internal/* ]] || [[ "$FILE_PATH" == docs/architecture.md ]]; then
    ADR_COUNT=$(find docs/decisions -name 'ADR-*.md' 2>/dev/null | wc -l)
    if [ "$ADR_COUNT" -eq 0 ]; then
        echo "[doc-sync] No ADRs found. Record architectural decisions in docs/decisions/."
    fi
fi
