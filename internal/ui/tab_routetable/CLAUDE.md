# Tab Route Table

## Role
Implements the Route Tables tab: lists route tables with inline search, and shows detail overlays for individual route entries or associated subnet IDs.

## Key Files
- `model.go` — `RouteTableModel` (implements `TabModel`); action menu, search model, update handlers, `loadRouteTables`
- `table.go` — Column definitions (name, RT ID, VPC, main flag, subnet count, route count), `RenderTable`, status bar
- `detail.go` — `RenderRouteEntries` (destination→target table), `RenderSubnets` (associated subnet list)

## Rules
- Detail overlays use locally-cached data (no async load); `detailKind` enum selects `detailRouteEntries` vs `detailSubnets`
- Action menu offers: Route Entries and Associated Subnets — no cross-tab navigation
- Search filters on name and route table ID only
- `RenderSubnets` notes when no explicit associations exist and the table is the main RT for the VPC
- Detail overlay closes on Esc (not any key, unlike VPC/Subnet tabs)
- Compact columns activate below 80-char terminal width
