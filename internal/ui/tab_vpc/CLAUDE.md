# Tab VPC

## Role
Implements the VPC tab: lists VPCs with inline search, and shows a rich detail overlay of all sub-resources (IGWs, NAT gateways, peering connections, TGW attachments, VPC endpoints, Elastic IPs).

## Key Files
- `model.go` — `VPCModel` (implements `TabModel`); action menu, search model, update handlers, `loadVPCs`
- `table.go` — Column definitions, `RenderTable`, `cellValue`/`cellStyle`, status bar renderer
- `detail.go` — `vpcDetailData`, `loadVPCDetail` (fetches all sub-resources filtered to VPC), `RenderVPCDetail`

## Rules
- Detail data is fetched lazily on "detail" action; `detailLoading` flag gates rendering until `vpcDetailLoadedMsg` arrives
- Action menu offers: VPC Details (loads sub-resources), and cross-tab navigation to Subnets, Routes, SG tabs via `shared.NavigateToTab`
- Search filters on name, VPC ID, and CIDR block
- `vpcDetailData` filters each sub-resource slice down to the selected VPC before storing
- Compact columns activate below 80-char terminal width; state column uses `StateRunning`/`StatePending` styles
- Any key press while `vsDetail` (not loading) closes the overlay and returns to `vsTable`
