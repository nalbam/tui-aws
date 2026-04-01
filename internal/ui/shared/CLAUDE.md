# Shared UI

## Role
Foundation package providing the TabModel interface, SharedState, reusable table/overlay primitives, and all Gruvbox lipgloss styles used across every tab.

## Key Files
- `tab.go` — `TabModel` interface, `TabID` enum, `SharedState` struct, `NavigateToTab` message, `CachedData` helpers
- `styles.go` — All lipgloss style vars (status bar, tabs, table, state colors, overlays) + `StateStyle(state)` helper
- `table.go` — `Column` type, `ExpandNameColumn` (fills terminal width), `RenderRow` (truncate + pad + style)
- `overlay.go` — `RenderOverlay` (rounded border), `PlaceOverlay` (centers on screen, accounts for tab/help bars)
- `selector.go` — `SelectorModel` generic list picker used by profile/region selectors

## Rules
- All tabs import `shared`; `shared` must never import any tab package (no circular deps)
- `SharedState` is passed by pointer to `Init`/`Update`/`View` — tabs mutate profile/region/cache via it
- Cache keys are `prefix::profile::region`; use `GetCache`/`SetCache`/`ClearCache` methods
- `ExpandNameColumn` clamps the `name` column to 20–60 chars; always call it before rendering
- `RenderRow` applies ANSI truncation then padding before styling — style must come last
- Gruvbox color palette is the sole color scheme; do not introduce new hex colors outside this file
