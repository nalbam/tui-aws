# ELB Tab
## Role
Load Balancers (ALB/NLB/CLB) — list LBs with interactive target group detail.
## Key Files
- `model.go` — ELBModel with vsDetail (interactive TG cursor) and vsTGDetail states
- `table.go` — Type column: ALB=blue, NLB=green, GWLB=yellow, CLB=gray
- `detail.go` — Interactive detail: SGs, listeners, selectable TG list → target health
## Rules
- Uses `elbv2.Client` (ALB/NLB/GWLB) and `elb.Client` (Classic) in parallel
- Two separate AWS clients — elbv2 for modern LBs, elb for Classic LBs
- Target groups selectable with ↑↓ cursor; Enter shows registered targets with health status
- Target health: healthy=green, unhealthy=red, draining/initial=yellow
- Detail has two sub-states: vsDetail (LB info + TG list) and vsTGDetail (target health view)
- CLB doesn't have target groups — detail shows listeners and instances instead
- Action menu offers: "LB Details" (loads listeners, TGs, and SGs)
- Search filters on LB name and DNS name
- Overlay closes on Esc; Esc from vsTGDetail returns to vsDetail
