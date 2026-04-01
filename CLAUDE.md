# Project Context

## Overview
tui-aws — AWS EC2 인스턴스 관리, VPC 네트워크 인프라 탐색, 연결성 검사를 위한 Go TUI 도구.

## Tech Stack
- **Language:** Go 1.25
- **TUI:** Bubble Tea v2 (Elm architecture), Lip Gloss v2 (Gruvbox theme)
- **AWS:** aws-sdk-go-v2 (EC2, SSM, STS) — single ec2.Client for all VPC/Subnet/SG/NACL/RT APIs
- **SSM Session:** `os/exec` → `aws ssm start-session` via custom `ssmExecCmd` + `tea.Exec()`
- **Build:** Makefile, cross-compile (linux/darwin × amd64/arm64)

## Project Structure
```
main.go                          Entry point, config migration, TUI launch
scripts/setup.sh                 Cross-platform setup & install script
internal/
  config/                        Config load/save (~/.tui-aws/config.json)
  store/                         Favorites & history CRUD (~/.tui-aws/)
  aws/
    ec2.go                       Instance model, FetchInstances, EnrichVpcSubnetInfo
    vpc.go                       VPC, IGW, NAT, Peering, TGW, Endpoint, EIP
    subnet.go                    Subnet, ENI
    network.go                   RouteTable, Route entries
    security.go                  SecurityGroup rules, NetworkACL rules
    reachability.go              VPC Reachability Analyzer
    profile.go                   AWS profile parsing (~/.aws/credentials + config)
    session.go                   SDK client factory (EC2/SSM/STS)
    ssm.go                       SSM command building, prerequisite checks
  ui/
    root.go                      RootModel (tea.Model), tab switching, global overlays
    tab.go                       Re-exports: TabModel, TabID, SharedState, NavigateToTab
    placeholder.go               PlaceholderTab for future tabs
    shared/
      tab.go                     TabModel interface, SharedState, CachedData, TabID enum
      styles.go                  All Lip Gloss styles (Gruvbox), tab bar styles
      table.go                   Column, RenderRow, ExpandNameColumn
      overlay.go                 RenderOverlay, PlaceOverlay (centered)
      selector.go                SelectorModel (generic list picker)
    tab_ec2/                     EC2 tab: SSM, port forward, Network Path, favorites
    tab_vpc/                     VPC tab: list + details (IGW/NAT/Peering/TGW/Endpoint/EIP)
    tab_subnet/                  Subnet tab: list + ENI viewer
    tab_routetable/              Route Table tab: list + route entries
    tab_sg/                      SG/NACL tab: rules viewer (f toggles mode)
    tab_troubleshoot/            Connectivity checker + Reachability Analyzer
docs/                            Architecture docs, ADRs, runbooks, specs
.claude/                         Claude settings, hooks, skills
```

## Conventions
- **Tab architecture:** RootModel owns SharedState, each tab implements TabModel interface
- **SharedState** in `shared/` package to avoid circular imports; `ui/tab.go` re-exports
- **EC2Model** sends `SSMExecRequest` messages; RootModel intercepts and runs `tea.Exec`
- **Lazy loading:** tabs fetch data on first switch, 30s cache TTL
- **ssmExecCmd:** wraps exec.Cmd with `stty sane` + stdin TCIFLUSH after SSM session
- **InterruptFilter:** blocks OS SIGINT (raw mode delivers Ctrl+C as KeyPressMsg)
- **View()** always sets `v.AltScreen = true` (Bubble Tea v2 API)
- **Cell-width aware:** `lipgloss.Width()` + `ansi.Truncate()` for Unicode/emoji columns
- **ExpandNameColumn:** Name column fills remaining terminal width (min 20, max 60)
- Test files colocated: `*_test.go` alongside implementation
- JSON config/store files under `~/.tui-aws/`

## Key Commands
```bash
make build          # Build binary (tui-aws)
make build-all      # Cross-compile for all platforms
make test           # Run all tests (go test ./... -v)
make clean          # Remove build artifacts
go vet ./...        # Static analysis
go test ./internal/ui/tab_troubleshoot/ -v  # Connectivity checker tests
./scripts/setup.sh  # Install prerequisites + build
```

---

## Auto-Sync Rules

Rules below are applied automatically after Plan mode exit and on major code changes.

### Post-Plan Mode Actions
After exiting Plan mode (`/plan`), before starting implementation:

1. **Architecture decision made** -> Update `docs/architecture.md`
2. **Technical choice/trade-off made** -> Create `docs/decisions/ADR-NNN-title.md`
3. **New tab added** -> Create `tab_<name>/` package with model.go, table.go, detail.go
4. **New module added** -> Create `CLAUDE.md` in that module directory
5. **Operational procedure defined** -> Create runbook in `docs/runbooks/`
6. **Changes needed in this file** -> Update relevant sections above

### Code Change Sync Rules
- New directory under `internal/` -> Must create `CLAUDE.md` alongside
- New AWS API usage -> Update `internal/aws/CLAUDE.md`
- New tab added -> Register in `root.go`, update `shared/tab.go` TabID enum
- UI shared component changed -> Update `internal/ui/shared/` CLAUDE.md or inline docs
- Config/store schema changed -> Update respective module `CLAUDE.md`
- Infrastructure changed -> Update `docs/architecture.md` Infrastructure section

### ADR Numbering
Find the highest number in `docs/decisions/ADR-*.md` and increment by 1.
Format: `ADR-NNN-concise-title.md`
