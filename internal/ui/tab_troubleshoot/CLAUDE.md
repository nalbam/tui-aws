# Tab Troubleshoot

## Role
Implements the Connectivity Check tab: a multi-step form that lets users pick source/destination EC2 instances, then performs local SG+NACL+route validation and optionally runs AWS Reachability Analyzer.

## Key Files
- `model.go` — `TroubleshootModel` (implements `TabModel`); 6 view states (form, picker, result, confirmRA, raRunning, raResult), all render and update handlers
- `checker.go` — `CheckConnectivity` (pure function: 5-step check), plus `cidrContains`, `portMatches`, `protocolMatches` helpers
- `result.go` — `RenderResult` formats the step-by-step check output with pass/fail/skip icons
- `checker_test.go` — Unit tests for helper functions and full connectivity scenarios including cross-VPC with TGW

## Rules
- `CheckConnectivity` is a pure function — no AWS calls; all data must be pre-fetched and passed in
- Check order: Source SG Outbound → Source NACL Outbound → Source Route → Dest NACL Inbound → Dest SG Inbound
- Once a step fails, all subsequent steps are marked `Skipped: true`
- NACL rules are evaluated in ascending rule-number order; first match (allow or deny) wins
- Route lookup: explicit subnet association first, fallback to VPC's main route table
- `cidrContains` silently returns false for non-CIDR sources (sg-xxx, pl-xxx) — SG-reference rules are not evaluated
- Reachability Analyzer requires user confirmation (warns about potential AWS costs) before calling the API
- `IsEditing()` is exported for RootModel to suppress global key handling when protocol/port fields are active
