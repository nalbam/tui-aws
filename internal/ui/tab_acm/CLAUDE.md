# ACM Tab
## Role
ACM Certificates — list certs with domain, status, expiry, SANs, in-use resources.
## Key Files
- `model.go` — ACMModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Status: ISSUED=green, PENDING=yellow, EXPIRED/REVOKED/FAILED=red
- `detail.go` — Full cert detail + SANs list + InUseBy resource ARNs
## Rules
- Uses `acm.Client` (ListCertificates + DescribeCertificate per cert for full details)
- Two-step fetch: list returns summary, then describe per cert fills SANs/InUseBy/expiry
- Expiry date highlighted when approaching
- Detail is rendered synchronously from cached cert data
- Action menu offers only "Certificate Details"
- Search filters on domain name and certificate ARN
- Overlay closes on Esc; standard 4-state viewState flow
