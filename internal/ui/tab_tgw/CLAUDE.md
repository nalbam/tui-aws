# TGW Tab
## Role
Transit Gateways — list TGWs with attachments, route tables, and routes.
## Key Files
- `model.go` — TGWModel with table + detail (attachments + routes loaded on demand)
- `table.go` — Columns: Name, TGW ID, State, ASN, Attachments count
- `detail.go` — Attachments list + route tables with route entries (lazy loaded)
## Rules
- Uses `ec2.Client` (DescribeTransitGateways, DescribeTransitGatewayAttachments, SearchTransitGatewayRoutes)
- Route tables and routes fetched on demand when viewing detail (not at list time)
- Detail is an async multi-step load: first attachments, then route tables per TGW
- SearchTransitGatewayRoutes requires a filter — uses type=static+propagated
- Action menu offers: "TGW Details" (shows attachments + route tables)
- Search filters on TGW ID and name tag
- Overlay closes on Esc; standard viewState flow
