# EBS Tab
## Role
EBS Volumes — list volumes with state, type, size, IOPS, encryption status, attachments.
## Key Files
- `model.go` — EBSModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Columns include Encrypted (green check / red cross color-coded)
- `detail.go` — Full volume detail: type, size, IOPS, throughput, encryption, attachment info (instance ID + device)
## Rules
- Uses existing `ec2.Client` (DescribeVolumes)
- Encryption column styled green (✓) / red (✗) for visual compliance check
- Volume state uses standard state colors (in-use=green, available=yellow, etc.)
- Attachment info shows instance ID and device name
- Action menu offers only "Volume Details"
- Search filters on volume ID and name tag
- Overlay closes on Esc; standard 4-state viewState flow
