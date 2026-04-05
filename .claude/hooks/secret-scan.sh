#!/bin/bash
# Pre-commit hook: scan staged files for potential secrets.
# Blocks commit if secrets are detected.

cd "$(git rev-parse --show-toplevel 2>/dev/null)" || exit 0

STAGED_FILES=$(git diff --cached --name-only --diff-filter=ACM 2>/dev/null)
[ -z "$STAGED_FILES" ] && exit 0

ISSUES=0

for FILE in $STAGED_FILES; do
    [ ! -f "$FILE" ] && continue

    # Skip binary files and known safe patterns
    case "$FILE" in
        *.png|*.jpg|*.gif|*.ico|*.woff|*.ttf|*.eot|go.sum|*.lock) continue ;;
    esac

    # Check for common secret patterns
    if grep -nEi '(AKIA[0-9A-Z]{16}|password\s*[:=]\s*["\x27][^"\x27]{8,}|aws_secret_access_key|private_key|BEGIN (RSA |EC |DSA |OPENSSH )?PRIVATE KEY)' "$FILE" 2>/dev/null; then
        echo "[secret-scan] Potential secret found in $FILE"
        ISSUES=$((ISSUES + 1))
    fi

    # Check for .env-style key=value with long values (potential tokens)
    if [[ "$FILE" == *.env* ]] || [[ "$FILE" == *credentials* ]] || [[ "$FILE" == *secret* ]]; then
        echo "[secret-scan] Sensitive file staged: $FILE"
        ISSUES=$((ISSUES + 1))
    fi
done

if [ "$ISSUES" -gt 0 ]; then
    echo "[secret-scan] BLOCKED: $ISSUES potential secret(s) detected. Review before committing."
    exit 1
fi

exit 0
