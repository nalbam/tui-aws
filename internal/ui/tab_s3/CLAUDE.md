# S3 Tab
## Role
S3 Buckets — list all buckets (global), view versioning/encryption/public access on demand.
## Key Files
- `model.go` — S3Model implementing TabModel, detail loads on demand
- `table.go` — Columns: Name, Region, Created, Versioning, Encryption, Public
- `detail.go` — Bucket detail with versioning/encryption/public access block settings
## Rules
- Uses `s3.Client` (ListBuckets is global — returns all regions)
- S3 is a global service — ListBuckets returns buckets across all regions
- Detail fetches 3 separate APIs: GetBucketVersioning, GetBucketEncryption, GetPublicAccessBlock
- Detail loading is asynchronous (per-bucket API calls on demand)
- Public access block settings shown as individual boolean flags
- Action menu offers only "Bucket Details"
- Search filters on bucket name and region
- Overlay closes on Esc; standard viewState flow
