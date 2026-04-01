# Tab Subnet

## Role
Implements the Subnet tab: lists subnets with inline search, and shows an ENI detail overlay listing all elastic network interfaces attached to the selected subnet.

## Key Files
- `model.go` — `SubnetModel` (implements `TabModel`); action menu, search model, update handlers, `loadSubnets`
- `table.go` — Column definitions (name, subnet ID, VPC, CIDR, AZ, available IPs, public flag), `RenderTable`, status bar
- `detail.go` — `loadENIs` (fetches all ENIs filtered to subnet), `RenderENIDetail`

## Rules
- ENIs are fetched lazily on "enis" action; `detailLoading` gates rendering until `eniLoadedMsg` arrives
- Action menu offers: ENIs in this Subnet, and cross-tab navigation to VPC tab via `shared.NavigateToTab`
- Search filters on name, subnet ID, VPC ID, CIDR block, and AZ
- AZ column displays only the last 2 characters (e.g. "1a") to save space
- Public subnet column uses `StatePending` (yellow) style to draw attention
- Compact columns activate below 100-char terminal width
- Any key press while `vsDetail` (not loading) closes the overlay and returns to `vsTable`
