# CloudFront Tab
## Role
CloudFront Distributions — list distributions with origins, aliases, WAF, certificates.
## Key Files
- `model.go` — CFModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Columns: ID, Domain, Status, Enabled, Origins, Aliases
- `detail.go` — Full distribution detail + origins list + aliases + WAF Web ACL + certificate ARN
## Rules
- Uses `cloudfront.Client` (ListDistributions)
- CloudFront is a global service — client uses configured region but data is global
- Detail is rendered synchronously (no async load) from cached distribution data
- Action menu offers only "Distribution Details"
- Search filters on distribution ID, domain name, and aliases
- Overlay closes on Esc; standard 4-state viewState flow
