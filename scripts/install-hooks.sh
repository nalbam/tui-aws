#!/bin/bash
# Install git hooks for the tui-aws project.
# Run once after cloning: bash scripts/install-hooks.sh

set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel 2>/dev/null)
if [ -z "$REPO_ROOT" ]; then
    echo "Error: not inside a git repository"
    exit 1
fi

HOOKS_DIR="$REPO_ROOT/.git/hooks"
mkdir -p "$HOOKS_DIR"

# Install commit-msg hook: removes Co-Authored-By lines from commit messages
cat > "$HOOKS_DIR/commit-msg" << 'HOOK_EOF'
#!/bin/bash
# Remove Co-Authored-By lines (AI contributor exclusion)
COMMIT_MSG_FILE="$1"
if [ -f "$COMMIT_MSG_FILE" ]; then
    sed -i.bak '/^Co-[Aa]uthored-[Bb]y:/d' "$COMMIT_MSG_FILE"
    rm -f "${COMMIT_MSG_FILE}.bak"
fi
HOOK_EOF
chmod +x "$HOOKS_DIR/commit-msg"
echo "Installed: commit-msg (AI contributor exclusion)"

# Install pre-commit hook: secret scanning
cat > "$HOOKS_DIR/pre-commit" << 'HOOK_EOF'
#!/bin/bash
# Pre-commit: run secret scan
if [ -f .claude/hooks/secret-scan.sh ]; then
    bash .claude/hooks/secret-scan.sh
fi
HOOK_EOF
chmod +x "$HOOKS_DIR/pre-commit"
echo "Installed: pre-commit (secret scanning)"

echo ""
echo "Git hooks installed successfully."
echo "  commit-msg  — removes Co-Authored-By lines"
echo "  pre-commit  — scans for secrets before commit"
