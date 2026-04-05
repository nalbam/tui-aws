# R53 Tab
## Role
Route 53 Hosted Zones — list zones, view DNS records on demand.
## Key Files
- `model.go` — R53Model with zone list + record detail (lazy loaded)
- `table.go` — Columns: Name, ID, Private, Records count, Comment
- `detail.go` — Zone info + scrollable record list (Name, Type, TTL, Value/Alias)
## Rules
- Uses `route53.Client` (ListHostedZones, ListResourceRecordSets)
- Records loaded on demand when user opens zone detail — not fetched at list time
- Route 53 is a global service — data is region-independent
- Record values: standard records show value; alias records show alias target
- Detail supports scrolling for zones with many records
- Action menu offers only "Zone Details" (triggers record fetch)
- Search filters on zone name and hosted zone ID
- Overlay closes on Esc; standard viewState flow
