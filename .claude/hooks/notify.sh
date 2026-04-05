#!/bin/bash
# Notification hook: send webhook on important events.
# Configure NOTIFY_WEBHOOK_URL environment variable to enable.

[ -z "$NOTIFY_WEBHOOK_URL" ] && exit 0

EVENT="${1:-unknown}"
MESSAGE="${2:-Claude Code notification}"

# Send webhook (non-blocking, best-effort)
curl -s -X POST "$NOTIFY_WEBHOOK_URL" \
    -H "Content-Type: application/json" \
    -d "{\"event\": \"$EVENT\", \"message\": \"$MESSAGE\", \"project\": \"tui-aws\", \"timestamp\": \"$(date -u +%Y-%m-%dT%H:%M:%SZ)\"}" \
    --max-time 5 \
    >/dev/null 2>&1 &

exit 0
