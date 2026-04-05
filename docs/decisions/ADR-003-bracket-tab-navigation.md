# ADR-003: Bracket key tab navigation for 22 tabs

## Status
Accepted

## Context
The original tab navigation used number keys (1-9) to switch tabs directly. When the tab count grew from 6 to 22, number keys could only address the first 9 tabs. We needed a navigation scheme that scales to any number of tabs.

## Decision
Replace number key navigation with `[` / `]` bracket keys for prev/next tab cycling. Also support `Tab` / `Shift+Tab` as alternative bindings. The Check tab (tab_troubleshoot) suppresses these keys when `IsEditing()` is true (text input fields for protocol/port).

## Consequences

### Positive
- Scales to any number of tabs without key conflicts
- Intuitive left/right navigation matches visual tab bar layout
- Tab/Shift+Tab is a familiar convention from other TUI tools

### Negative
- No direct jump to a specific tab (must cycle through)
- Slightly slower to reach distant tabs (mitigated: 22 tabs wraps around)

### Risks
- None — the `[`/`]` keys don't conflict with any tab-level keybindings
