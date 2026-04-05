# Runbook: Adding a New Tab

## Purpose
Step-by-step guide for adding a new AWS service tab to tui-aws.

## Prerequisites
- [ ] Go 1.25+ installed
- [ ] Familiarity with Bubble Tea v2 Elm architecture (Model-View-Update)
- [ ] AWS service SDK imported in go.mod

## Procedure

### Step 1: Add TabID
Edit `internal/ui/shared/tab.go`:
- Add `Tab<Name>` to the TabID enum (before the closing parenthesis)
- Add a case in `Label()` returning the tab bar display name

### Step 2: Add AWS client (if needed)
Edit `internal/aws/session.go`:
- Add the new service client field to `Clients` struct
- Initialize it in `NewClients()`

### Step 3: Create AWS data functions
Create `internal/aws/<service>.go`:
- Define data structs (e.g., `type MyResource struct { ... }`)
- Implement `Fetch<Resources>(ctx, clients)` function with pagination

### Step 4: Create tab package
Create `internal/ui/tab_<name>/` with 3 files:
```bash
mkdir internal/ui/tab_<name>
```
- `model.go` — `<Name>Model` implementing `TabModel`, 4-state viewState, async load
- `table.go` — Column definitions, `RenderTable`, status bar
- `detail.go` — Detail overlay render function

### Step 5: Register tab
Edit `internal/ui/root.go`:
- Import the new tab package
- Add `tab_<name>.New<Name>Model()` to the tabs slice

### Step 6: Create documentation
- `internal/ui/tab_<name>/CLAUDE.md` — Role, key files, rules
- Update `internal/aws/CLAUDE.md` — Add new service to the list

### Step 7: Build and test
```bash
make build && ./tui-aws
go vet ./...
```

## Verification
- New tab appears in the tab bar
- `[`/`]` keys cycle through including the new tab
- Data loads on first tab activation (lazy loading)
- Detail overlay opens and closes cleanly

## Related
- ADR: docs/decisions/ADR-003-bracket-tab-navigation.md
- Architecture: docs/architecture.md#tab-architecture
