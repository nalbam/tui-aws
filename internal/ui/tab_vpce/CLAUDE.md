# VPCE Tab
## Role
VPC Endpoints — list endpoints with service name, type, state, and full details.
## Key Files
- `model.go` — VPCEModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Columns: Name, Endpoint ID, Service Name, Type, VPC, State
- `detail.go` — Full endpoint detail: subnets, route tables, SGs, ENIs, private DNS, creation time
## Rules
- Uses existing `ec2.Client` (DescribeVpcEndpoints)
- Type: Gateway (route table based) vs Interface (ENI based) — different detail sections
- Gateway endpoints show route tables; Interface endpoints show subnets, SGs, ENIs
- Private DNS name shown when available
- Action menu offers only "Endpoint Details"
- Search filters on endpoint ID, service name, and VPC ID
- Overlay closes on Esc; standard 4-state viewState flow
