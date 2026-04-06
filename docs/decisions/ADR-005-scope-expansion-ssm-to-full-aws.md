# ADR-005: Scope expansion from SSM tool to full AWS infrastructure TUI

## Status
Accepted

## Context
The project originally started as `tui-ssm`, a focused tool for managing SSM sessions to EC2 instances. As usage grew, switching between the TUI and AWS Console for networking, load balancing, DNS, and container resources became a constant friction point. Users needed one terminal-based tool that could explore the full AWS infrastructure, not just EC2/SSM.

The key question was whether to keep a focused SSM tool or expand into a comprehensive AWS infrastructure TUI covering VPC networking, containers, serverless, CDN, DNS, security, and monitoring.

## Decision
Rename the project from `tui-ssm` to `tui-aws` and adopt a phased expansion strategy:

1. **Phase 0** — Refactor from single-view to tab-based architecture (RootModel + TabModel interface)
2. **Phase 1** — VPC, Subnet tabs with detail overlays
3. **Phase 2** — Routes, SG/NACL tabs with Network Path visualization
4. **Phase 3** — Connectivity checker and Reachability Analyzer
5. **Subsequent phases** — ELB, ASG, EBS, TGW, CloudFront, WAF, ACM, R53, RDS, S3, ECS, EKS, Lambda, CloudWatch, IAM

Config and store paths migrated from `~/.tui-ssm/` to `~/.tui-aws/`. Module path changed from `tui-ssm` to `tui-aws`.

Relevant commits: `ca2d8cb` (module rename), `525f0d9` (phase 0 refactor), `8f9a4f3` (phase 1), `4bdcca3` (phase 2), `a7235f2` (phase 3).

## Consequences

### Positive
- Single tool replaces multiple Console tabs and CLI commands
- Tab-based architecture scales cleanly — each tab is an independent package
- SharedState (profile, region, cache, clients) is reused across all 22 tabs
- Cross-resource navigation (EC2 -> VPC -> Subnet -> SG) is natural in a unified TUI

### Negative
- Binary size grows with each AWS SDK service client (~25MB currently)
- More AWS API permissions required (though read-only for most tabs)
- Larger surface area to maintain (107 source files, ~22K lines)

### Risks
- Feature creep — new AWS services constantly launch (mitigated: only add tabs with clear user value)
- SDK dependency bloat — 18 AWS service clients in go.mod (mitigated: each client is imported only once in session.go)
