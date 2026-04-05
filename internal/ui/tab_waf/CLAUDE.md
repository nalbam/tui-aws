# WAF Tab
## Role
WAFv2 Web ACLs — list ACLs with rule count, default action, associated resources.
## Key Files
- `model.go` — WAFModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Default action: Allow=green, Block=red
- `detail.go` — Full ACL detail: rules list + associated resource ARNs
## Rules
- Uses `wafv2.Client` with REGIONAL scope only — CLOUDFRONT scope requires us-east-1 region
- Fetches rules count via `GetWebACL`, resources via `ListResourcesForWebACL`
- Two-step fetch: list returns summaries, then per-ACL describe fills rules and resources
- Default action color-coded: Allow=green, Block=red
- Action menu offers only "ACL Details"
- Search filters on Web ACL name and ID
- Overlay closes on Esc; standard 4-state viewState flow
