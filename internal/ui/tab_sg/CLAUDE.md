# Tab SG

## Role
Implements the Security Groups / Network ACLs tab: a single tab that toggles between SG and NACL views, with inline search and detail overlays for inbound/outbound rules.

## Key Files
- `model.go` — `SGModel` (implements `TabModel`); dual-mode state (`mode`: "sg"/"nacl"), lazy load with `sgLoaded`/`naclLoaded` flags, update handlers
- `table.go` — Column definitions and render functions for both SG (`RenderSGTable`) and NACL (`RenderNACLTable`) tables, status bar
- `detail.go` — `RenderSGRules` and `RenderNACLRules` (inbound or outbound overlay based on `detailKind`)

## Rules
- `f` key toggles between SG and NACL mode; each mode's data is fetched once and cached (`sgLoaded`/`naclLoaded`)
- Action menu for both SGs and NACLs offers only Inbound Rules and Outbound Rules
- Selected item for detail is snapshot-copied into `selectedSG`/`selectedNACL` at action-menu confirmation time
- NACL rule numbers are rendered as `*` when rule number is 32767 (the implicit deny-all)
- Detail overlay closes on Esc; `detailKind` (`detailInbound`/`detailOutbound`) is set in `updateActionMenu`
- SG compact columns below 100-char; NACL compact columns below 80-char terminal width
- Search filters on name and ID for whichever mode is active
