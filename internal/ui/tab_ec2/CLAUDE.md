# Tab EC2

## Role
Implements the EC2 instances tab: lists instances with search/filter/sort, opens SSM sessions and port-forward tunnels, and shows detail overlays for SGs, instance info, and full network path.

## Key Files
- `model.go` — `EC2Model` (implements `TabModel`); all `Update` dispatch, SSM/port-forward cmd helpers, `loadNetworkPath`
- `table.go` — Column definitions, `RenderTable`, `SortInstances` (favs > recent > field), `FilterBySearch/State`
- `actions.go` — `ActionMenuModel`, `PortForwardModel`, `RenderSecurityGroups`, `RenderInstanceDetail`, `RenderNetworkPath`
- `search.go` — `SearchModel` (Insert/Backspace/Clear/Render)
- `filter.go` — `FilterModel` state-checkbox picker with `Toggle`/`ClearAll`/`Label`/`Render`
- `table_test.go` — Unit tests for sort priority and search/state filter logic

## Rules
- `SSMExecRequest` is a message sent up to RootModel — never call `tea.Exec` directly from this tab
- `SSMSessionDoneMsg` is exported; RootModel matches it to record history then forwards it back here
- Sort priority: favorites first, recent history second, user-selected column third
- `FilterBySearch` matches on name, instance ID, and private IP (case-insensitive)
- `loadNetworkPath` falls back to the VPC's main route table when no explicit subnet association exists
- SG matching in `loadNetworkPath` uses SG name (not ID) against `inst.SecurityGroups`
- Detail overlay (`showDetail`) stays active while `actionMenu.Active` is true; any key dismisses it
- Compact columns activate below 80-char terminal width
