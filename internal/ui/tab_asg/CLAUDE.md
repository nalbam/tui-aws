# ASG Tab
## Role
Auto Scaling Groups — list groups with min/max/desired, instances, scaling policies, target groups.
## Key Files
- `model.go` — ASGModel implementing TabModel, 4-state viewState (table/search/actionMenu/detail)
- `table.go` — Columns: Name, Min, Max, Desired, Instances, Health, AZs
- `detail.go` — Full ASG detail: launch config/template, instances list, target groups, scaling policies
## Rules
- Uses `autoscaling.Client` (DescribeAutoScalingGroups) — separate from EC2 client
- Instance count shown as running/total
- Launch config vs launch template: displays whichever is non-empty
- Detail lists instance IDs, target group ARNs, and scaling policies with type/adjustments
- AZs displayed as comma-separated list
- Action menu offers only "ASG Details"
- Search filters on ASG name
- Overlay closes on Esc; standard 4-state viewState flow
