# RDS Tab
## Role
RDS DB Instances — list instances with engine, class, endpoint, multi-AZ, storage.
## Key Files
- `model.go` — RDSModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Status: available=green, creating/modifying=yellow, deleting/failed=red
- `detail.go` — Full DB detail: endpoint, port, VPC, subnet group, SGs, encryption, public access
## Rules
- Uses `rds.Client` (DescribeDBInstances)
- Detail includes security groups list and subnet group name
- MultiAZ displayed as boolean label
- Storage shows type + allocated GiB
- Action menu offers only "Instance Details"
- Search filters on DB identifier and engine type
- Overlay closes on Esc; standard 4-state viewState flow
